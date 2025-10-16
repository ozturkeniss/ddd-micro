package http

import (
	"net/http"
	"strings"

	"github.com/ddd-micro/internal/product/application"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication and authorization for HTTP
type AuthMiddleware struct {
	userService *application.UserService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userService *application.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

// AuthRequired middleware for endpoints that require authentication
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token with user service
		user, err := m.userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User account is inactive",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", user.Id)
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		c.Set("user_active", user.IsActive)

		c.Next()
	}
}

// AdminRequired middleware for endpoints that require admin access
func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check authentication
		m.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		// Check if user is admin
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware for endpoints that can work with or without authentication
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>" format
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.Next()
			return
		}

		// Try to validate token, but don't fail if invalid
		user, err := m.userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}

		// Add user info to context if valid
		if user.IsActive {
			c.Set("user_id", user.Id)
			c.Set("user_email", user.Email)
			c.Set("user_role", user.Role)
			c.Set("user_active", user.IsActive)
		}

		c.Next()
	}
}
