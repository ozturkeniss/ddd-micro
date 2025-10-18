package http

import (
	"net/http"
	"strings"

	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication and authorization
type AuthMiddleware struct {
	userClient client.UserClient
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userClient client.UserClient) *AuthMiddleware {
	return &AuthMiddleware{
		userClient: userClient,
	}
}

// AuthRequired middleware for endpoints that require authentication
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Bearer token required",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token with user service
		user, err := m.userClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User account is inactive",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", uint(user.Id))
		c.Set("userRole", user.Role)
		c.Set("userEmail", user.Email)

		c.Next()
	}
}

// AdminRequired middleware for endpoints that require admin role
func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		m.AuthRequired()(c)

		// Check if request was aborted by auth middleware
		if c.IsAborted() {
			return
		}

		// Get user role from context
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Admin role required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
