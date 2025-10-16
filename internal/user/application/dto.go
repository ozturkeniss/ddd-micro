package application

import (
	"time"

	"github.com/ddd-micro/internal/user/domain"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserRequest represents the request to update a user (self-update)
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUserByAdminRequest represents the request for admin to update any user
type UpdateUserByAdminRequest struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      domain.Role `json:"role"`
	IsActive  *bool       `json:"is_active"`
}

// AssignRoleRequest represents the request to assign a role to a user
type AssignRoleRequest struct {
	Role domain.Role `json:"role" binding:"required"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// RefreshTokenRequest represents the request to refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenResponse represents the response with token
type TokenResponse struct {
	Token string `json:"token"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID        uint        `json:"id"`
	Email     string      `json:"email"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Role      domain.Role `json:"role"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// LoginResponse represents the login response with token
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// ListUsersResponse represents the paginated list of users
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"`
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
}

