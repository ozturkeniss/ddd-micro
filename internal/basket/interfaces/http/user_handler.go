package http

import (
	"net/http"

	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userClient client.UserClient
}

// NewUserHandler creates a new user handler
func NewUserHandler(userClient client.UserClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// GetProfile retrieves the user's profile information
// @Summary Get user profile
// @Description Retrieves the profile information for the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	user, err := h.userClient.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// ValidateToken validates a JWT token
// @Summary Validate token
// @Description Validates a JWT token and returns user information
// @Tags user
// @Accept json
// @Produce json
// @Param request body map[string]string true "Token validation request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/validate-token [post]
func (h *UserHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	user, err := h.userClient.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
