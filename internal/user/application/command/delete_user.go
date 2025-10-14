package command

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// DeleteUserCommand represents a command to delete a user
type DeleteUserCommand struct {
	UserID uint
}

// DeleteUserHandler handles the DeleteUserCommand
type DeleteUserHandler struct {
	repo domain.UserRepository
}

// NewDeleteUserHandler creates a new DeleteUserHandler
func NewDeleteUserHandler(repo domain.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{
		repo: repo,
	}
}

// Handle executes the DeleteUserCommand
func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	return h.repo.Delete(ctx, cmd.UserID)
}

