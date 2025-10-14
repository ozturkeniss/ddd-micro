package command

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// UpdateUserCommand represents a command to update a user (self-update)
type UpdateUserCommand struct {
	UserID    uint
	FirstName string
	LastName  string
}

// UpdateUserHandler handles the UpdateUserCommand
type UpdateUserHandler struct {
	repo domain.UserRepository
}

// NewUpdateUserHandler creates a new UpdateUserHandler
func NewUpdateUserHandler(repo domain.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{
		repo: repo,
	}
}

// Handle executes the UpdateUserCommand
func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error) {
	// Get existing user
	user, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if cmd.FirstName != "" {
		user.FirstName = cmd.FirstName
	}
	if cmd.LastName != "" {
		user.LastName = cmd.LastName
	}

	// Save changes
	if err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

