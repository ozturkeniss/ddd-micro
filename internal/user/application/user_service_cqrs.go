package application

import (
	"context"
	"time"

	"github.com/ddd-micro/internal/user/application/command"
	"github.com/ddd-micro/internal/user/application/query"
	"github.com/ddd-micro/internal/user/domain"
)

// UserServiceCQRS handles user business logic using CQRS pattern
type UserServiceCQRS struct {
	// Command handlers
	createUserHandler        *command.CreateUserHandler
	updateUserHandler        *command.UpdateUserHandler
	updateUserByAdminHandler *command.UpdateUserByAdminHandler
	deleteUserHandler        *command.DeleteUserHandler
	changePasswordHandler    *command.ChangePasswordHandler
	assignRoleHandler        *command.AssignRoleHandler
	loginHandler             *command.LoginHandler

	// Query handlers
	getUserByIDHandler    *query.GetUserByIDHandler
	getUserByEmailHandler *query.GetUserByEmailHandler
	listUsersHandler      *query.ListUsersHandler

	// JWT helper for token operations
	jwtHelper *JWTHelper
}

// NewUserServiceCQRS creates a new CQRS-based user service
func NewUserServiceCQRS(repo domain.UserRepository, jwtSecretKey string, tokenDuration time.Duration) *UserServiceCQRS {
	return &UserServiceCQRS{
		// Initialize command handlers
		createUserHandler:        command.NewCreateUserHandler(repo),
		updateUserHandler:        command.NewUpdateUserHandler(repo),
		updateUserByAdminHandler: command.NewUpdateUserByAdminHandler(repo),
		deleteUserHandler:        command.NewDeleteUserHandler(repo),
		changePasswordHandler:    command.NewChangePasswordHandler(repo),
		assignRoleHandler:        command.NewAssignRoleHandler(repo),
		loginHandler:             command.NewLoginHandler(repo, jwtSecretKey, tokenDuration),

		// Initialize query handlers
		getUserByIDHandler:    query.NewGetUserByIDHandler(repo),
		getUserByEmailHandler: query.NewGetUserByEmailHandler(repo),
		listUsersHandler:      query.NewListUsersHandler(repo),

		// Initialize JWT helper
		jwtHelper: NewJWTHelper(jwtSecretKey, tokenDuration),
	}
}

// ========== COMMAND METHODS ==========

// CreateUser creates a new user with hashed password
func (s *UserServiceCQRS) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	cmd := command.CreateUserCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := s.createUserHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Login authenticates a user and returns a JWT token
func (s *UserServiceCQRS) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	cmd := command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := s.loginHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  *s.toUserResponse(result.User),
		Token: result.Token,
	}, nil
}

// UpdateUser updates a user's information (self-update)
func (s *UserServiceCQRS) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error) {
	cmd := command.UpdateUserCommand{
		UserID:    id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := s.updateUserHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// UpdateUserByAdmin updates any user's information (admin only)
func (s *UserServiceCQRS) UpdateUserByAdmin(ctx context.Context, id uint, req UpdateUserByAdminRequest) (*UserResponse, error) {
	cmd := command.UpdateUserByAdminCommand{
		UserID:    id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  req.IsActive,
	}

	user, err := s.updateUserByAdminHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// DeleteUser soft deletes a user
func (s *UserServiceCQRS) DeleteUser(ctx context.Context, id uint) error {
	cmd := command.DeleteUserCommand{
		UserID: id,
	}

	return s.deleteUserHandler.Handle(ctx, cmd)
}

// ChangePassword changes a user's password
func (s *UserServiceCQRS) ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error {
	cmd := command.ChangePasswordCommand{
		UserID:      id,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	return s.changePasswordHandler.Handle(ctx, cmd)
}

// AssignRole assigns a role to a user (admin only)
func (s *UserServiceCQRS) AssignRole(ctx context.Context, userID uint, role domain.Role) (*UserResponse, error) {
	cmd := command.AssignRoleCommand{
		UserID: userID,
		Role:   role,
	}

	user, err := s.assignRoleHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// ========== QUERY METHODS ==========

// GetUserByID retrieves a user by ID
func (s *UserServiceCQRS) GetUserByID(ctx context.Context, id uint) (*UserResponse, error) {
	q := query.GetUserByIDQuery{
		UserID: id,
	}

	user, err := s.getUserByIDHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// GetUserByEmail retrieves a user by email
func (s *UserServiceCQRS) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	q := query.GetUserByEmailQuery{
		Email: email,
	}

	user, err := s.getUserByEmailHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// ListUsers retrieves all users with pagination
func (s *UserServiceCQRS) ListUsers(ctx context.Context, offset, limit int) (*ListUsersResponse, error) {
	q := query.ListUsersQuery{
		Offset: offset,
		Limit:  limit,
	}

	result, err := s.listUsersHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	userResponses := make([]UserResponse, len(result.Users))
	for i, user := range result.Users {
		userResponses[i] = *s.toUserResponse(user)
	}

	return &ListUsersResponse{
		Users:  userResponses,
		Total:  result.Total,
		Offset: result.Offset,
		Limit:  result.Limit,
	}, nil
}

// ========== TOKEN METHODS ==========

// ValidateToken validates a JWT token
func (s *UserServiceCQRS) ValidateToken(tokenString string) (*JWTClaims, error) {
	return s.jwtHelper.ValidateToken(tokenString)
}

// RefreshToken refreshes a JWT token
func (s *UserServiceCQRS) RefreshToken(tokenString string) (string, error) {
	return s.jwtHelper.RefreshToken(tokenString)
}

// ========== HELPER METHODS ==========

// toUserResponse converts domain.User to UserResponse
func (s *UserServiceCQRS) toUserResponse(user *domain.User) *UserResponse {
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
