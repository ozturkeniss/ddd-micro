package http

import (
	"strings"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
	UserRoleKey         = "user_role"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(userService *application.UserServiceCQRS) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			UnauthorizedResponse(c, "Token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := userService.ValidateToken(token)
		if err != nil {
			UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetUserRole retrieves user role from context
func GetUserRole(c *gin.Context) (domain.Role, bool) {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}
	return role.(domain.Role), true
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(requiredRole domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetUserRole(c)
		if !exists {
			UnauthorizedResponse(c, "User role not found")
			c.Abort()
			return
		}

		if role != requiredRole {
			ErrorResponse(c, 403, "Insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleAdmin)
}

// RequireRoles creates a middleware that requires one of the specified roles
func RequireRoles(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetUserRole(c)
		if !exists {
			UnauthorizedResponse(c, "User role not found")
			c.Abort()
			return
		}

		// Check if user has any of the allowed roles
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			ErrorResponse(c, 403, "Insufficient permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
