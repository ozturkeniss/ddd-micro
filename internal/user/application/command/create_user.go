package command

import (
	"context"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
)

// CreateUserCommand represents a command to create a new user
type CreateUserCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// CreateUserHandler handles the CreateUserCommand
type CreateUserHandler struct {
	repo           domain.UserRepository
	passwordHasher *application.PasswordHasher
}

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(repo domain.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{
		repo:           repo,
		passwordHasher: application.NewPasswordHasher(),
	}
}

// Handle executes the CreateUserCommand
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	// Hash the password
	hashedPassword, err := h.passwordHasher.HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := &domain.User{
		Email:     cmd.Email,
		Password:  hashedPassword,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Role:      domain.RoleUser,
		IsActive:  true,
	}

	// Save to repository
	if err := h.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

