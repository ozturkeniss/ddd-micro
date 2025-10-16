package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService interface{}
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService interface{}) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves the authenticated user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile (User only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// This would typically extract user ID from JWT token
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile endpoint - to be implemented",
	})
}

// ValidateToken validates a JWT token
// @Summary Validate JWT token
// @Description Validate a JWT token and return user information (User only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /users/validate-token [post]
func (h *UserHandler) ValidateToken(c *gin.Context) {
	// This would validate the JWT token
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Token validation endpoint - to be implemented",
	})
}
