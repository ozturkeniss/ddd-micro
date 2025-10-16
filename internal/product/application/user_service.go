package application

import (
	"context"
	"fmt"

	userpb "github.com/ddd-micro/api/proto/user"
	"github.com/ddd-micro/internal/product/infrastructure/client"
)

// UserService handles user-related operations for product service
type UserService struct {
	userClient client.UserClient
}

// NewUserService creates a new user service
func NewUserService(userClient client.UserClient) *UserService {
	return &UserService{
		userClient: userClient,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*userpb.User, error) {
	return s.userClient.GetUserByID(ctx, userID)
}

// ValidateToken validates a JWT token and returns user information
func (s *UserService) ValidateToken(ctx context.Context, token string) (*userpb.User, error) {
	return s.userClient.ValidateToken(ctx, token)
}

// IsUserAdmin checks if a user is admin
func (s *UserService) IsUserAdmin(ctx context.Context, userID uint) (bool, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	return user.Role == "admin", nil
}

// IsUserActive checks if a user is active
func (s *UserService) IsUserActive(ctx context.Context, userID uint) (bool, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	return user.IsActive, nil
}

// ValidateUserAccess validates if a user has access to perform operations
func (s *UserService) ValidateUserAccess(ctx context.Context, userID uint) error {
	// Check if user exists and is active
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return fmt.Errorf("user account is inactive")
	}

	return nil
}

// ValidateAdminAccess validates if a user has admin access
func (s *UserService) ValidateAdminAccess(ctx context.Context, userID uint) error {
	// Check if user exists and is active
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return fmt.Errorf("user account is inactive")
	}

	if user.Role != "admin" {
		return fmt.Errorf("admin access required")
	}

	return nil
}
