package application

import (
	"context"
	"errors"
	"time"

	"github.com/ddd-micro/internal/user/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

// UserService handles user business logic
type UserService struct {
	repo           domain.UserRepository
	passwordHasher *PasswordHasher
	jwtHelper      *JWTHelper
}

// NewUserService creates a new user service
func NewUserService(repo domain.UserRepository, jwtSecretKey string, tokenDuration time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		passwordHasher: NewPasswordHasher(),
		jwtHelper:      NewJWTHelper(jwtSecretKey, tokenDuration),
	}
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Hash the password
	hashedPassword, err := s.passwordHasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := &domain.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.RoleUser, // Default role is user
		IsActive:  true,
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := s.passwordHasher.ComparePassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtHelper.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  *s.toUserResponse(user),
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// UpdateUser updates a user's information (self-update)
func (s *UserService) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields (users can only update their name)
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// UpdateUserByAdmin updates any user's information (admin only)
func (s *UserService) UpdateUserByAdmin(ctx context.Context, id uint, req UpdateUserByAdminRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" && req.Role.IsValid() {
		user.AssignRole(req.Role)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// AssignRole assigns a role to a user (admin only)
func (s *UserService) AssignRole(ctx context.Context, userID uint, role domain.Role) (*UserResponse, error) {
	// Validate role
	if !role.IsValid() {
		return nil, errors.New("invalid role")
	}

	// Get user
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Assign role
	user.AssignRole(role)

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(ctx context.Context, offset, limit int) (*ListUsersResponse, error) {
	users, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(user)
	}

	return &ListUsersResponse{
		Users:  userResponses,
		Total:  len(userResponses),
		Offset: offset,
		Limit:  limit,
	}, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := s.passwordHasher.ComparePassword(user.Password, oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := s.passwordHasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	return s.repo.Update(ctx, user)
}

// ValidateToken validates a JWT token
func (s *UserService) ValidateToken(tokenString string) (*JWTClaims, error) {
	return s.jwtHelper.ValidateToken(tokenString)
}

// RefreshToken refreshes a JWT token
func (s *UserService) RefreshToken(tokenString string) (string, error) {
	return s.jwtHelper.RefreshToken(tokenString)
}

// toUserResponse converts domain.User to UserResponse
func (s *UserService) toUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
