package command

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

// LoginCommand represents a command to login
type LoginCommand struct {
	Email    string
	Password string
}

// LoginResult represents the result of login
type LoginResult struct {
	User  *domain.User
	Token string
}

// LoginHandler handles the LoginCommand
type LoginHandler struct {
	repo           domain.UserRepository
	passwordHasher *application.PasswordHasher
	jwtHelper      *application.JWTHelper
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(repo domain.UserRepository, jwtSecret string, tokenDuration int64) *LoginHandler {
	return &LoginHandler{
		repo:           repo,
		passwordHasher: application.NewPasswordHasher(),
		jwtHelper:      application.NewJWTHelper(jwtSecret, tokenDuration),
	}
}

// Handle executes the LoginCommand
func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// Get user by email
	user, err := h.repo.GetByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := h.passwordHasher.ComparePassword(user.Password, cmd.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := h.jwtHelper.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}

