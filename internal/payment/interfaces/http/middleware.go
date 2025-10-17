package http

import (
	"net/http"
	"strings"

	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(userClient client.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		user, err := userClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", uint(user.Id))
		c.Set("user_role", user.Role)
		c.Set("user_email", user.Email)

		c.Next()
	}
}

// RBACMiddleware checks if user has required permission
func RBACMiddleware(userClient client.UserClient, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		hasPermission, err := userClient.CheckPermission(c.Request.Context(), userID.(uint), resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnlyMiddleware ensures only admin users can access
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserOrAdminMiddleware allows both user and admin access
func UserOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if userRole != "user" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "User or admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
