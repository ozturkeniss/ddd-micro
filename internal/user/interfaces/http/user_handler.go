package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *application.UserServiceCQRS
	metrics     *monitoring.PrometheusMetrics
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *application.UserServiceCQRS, metrics *monitoring.PrometheusMetrics) *UserHandler {
	return &UserHandler{
		userService: userService,
		metrics:     metrics,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body application.CreateUserRequest true "User registration data"
// @Success 201 {object} Response{data=application.UserResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	// Start tracing span
	span, ctx := monitoring.StartSpanFromGinContext(c, "user.register")
	defer span.Finish()

	var req application.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanError(span, err)
		ValidationErrorResponse(c, err)
		return
	}

	start := time.Now()
	user, err := h.userService.CreateUser(ctx, req)
	duration := time.Since(start)

	// Record database query duration
	h.metrics.RecordDatabaseQuery("create_user", "users", duration)

	if err != nil {
		monitoring.LogSpanError(span, err)
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Record successful registration
	h.metrics.RecordUserRegistration()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"user.id":    user.ID,
		"user.email": user.Email,
		"success":    true,
	})

	SuccessResponse(c, http.StatusCreated, "User created successfully", user)
}

// Login handles user authentication
// @Summary Login user
// @Tags users
// @Accept json
// @Produce json
// @Param request body application.LoginRequest true "Login credentials"
// @Success 200 {object} Response{data=application.LoginResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	// Start tracing span
	span, ctx := monitoring.StartSpanFromGinContext(c, "user.login")
	defer span.Finish()

	var req application.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		monitoring.LogSpanError(span, err)
		ValidationErrorResponse(c, err)
		return
	}

	start := time.Now()
	loginResp, err := h.userService.Login(ctx, req)
	duration := time.Since(start)

	// Record database query duration
	h.metrics.RecordDatabaseQuery("user_login", "users", duration)

	if err != nil {
		monitoring.LogSpanError(span, err)
		if err == application.ErrInvalidCredentials || err == application.ErrUserInactive {
			h.metrics.RecordUserLoginFailure()
			UnauthorizedResponse(c, err.Error())
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Login failed", err)
		return
	}

	// Record successful login
	h.metrics.RecordUserLogin()
	monitoring.SetSpanTags(span, map[string]interface{}{
		"user.id":    loginResp.User.ID,
		"user.email": loginResp.User.Email,
		"success":    true,
	})

	SuccessResponse(c, http.StatusOK, "Login successful", loginResp)
}

// GetProfile retrieves the authenticated user's profile
// @Summary Get user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response{data=application.UserResponse}
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		NotFoundResponse(c, "User not found")
		return
	}

	SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// GetUserByID retrieves a user by ID
// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} Response{data=application.UserResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		NotFoundResponse(c, "User not found")
		return
	}

	SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// UpdateProfile updates the authenticated user's profile
// @Summary Update user profile
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body application.UpdateUserRequest true "Update data"
// @Success 200 {object} Response{data=application.UserResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req application.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

// DeleteUser deletes a user
// @Summary Delete user
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers retrieves all users with pagination
// @Summary List users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} Response{data=application.ListUsersResponse}
// @Failure 401 {object} Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Limit max page size
	if limit > 100 {
		limit = 100
	}

	users, err := h.userService.ListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

// ChangePassword handles password change
// @Summary Change password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body application.ChangePasswordRequest true "Password change data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req application.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err == application.ErrInvalidCredentials {
			UnauthorizedResponse(c, "Invalid old password")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to change password", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body application.RefreshTokenRequest true "Token to refresh"
// @Success 200 {object} Response{data=application.TokenResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /users/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	newToken, err := h.userService.RefreshToken(req.Token)
	if err != nil {
		UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	SuccessResponse(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"token": newToken,
	})
}

// ========== ADMIN HANDLERS ==========

// UpdateUserByAdmin updates any user's information (admin only)
// @Summary Update user by admin
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body application.UpdateUserByAdminRequest true "Update data"
// @Success 200 {object} Response{data=application.UserResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUserByAdmin(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req application.UpdateUserByAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.UpdateUserByAdmin(c.Request.Context(), uint(id), req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "User updated successfully by admin", user)
}

// AssignRole assigns a role to a user (admin only)
// @Summary Assign role to user
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body application.AssignRoleRequest true "Role data"
// @Success 200 {object} Response{data=application.UserResponse}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Router /admin/users/{id}/assign-role [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req application.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	user, err := h.userService.AssignRole(c.Request.Context(), uint(id), req.Role)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to assign role", err)
		return
	}

	SuccessResponse(c, http.StatusOK, "Role assigned successfully", user)
}
