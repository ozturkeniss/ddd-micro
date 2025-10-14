package command

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// UpdateUserByAdminCommand represents a command for admin to update any user
type UpdateUserByAdminCommand struct {
	UserID    uint
	FirstName string
	LastName  string
	Role      domain.Role
	IsActive  *bool
}

// UpdateUserByAdminHandler handles the UpdateUserByAdminCommand
type UpdateUserByAdminHandler struct {
	repo domain.UserRepository
}

// NewUpdateUserByAdminHandler creates a new UpdateUserByAdminHandler
func NewUpdateUserByAdminHandler(repo domain.UserRepository) *UpdateUserByAdminHandler {
	return &UpdateUserByAdminHandler{
		repo: repo,
	}
}

// Handle executes the UpdateUserByAdminCommand
func (h *UpdateUserByAdminHandler) Handle(ctx context.Context, cmd UpdateUserByAdminCommand) (*domain.User, error) {
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
	if cmd.Role != "" && cmd.Role.IsValid() {
		user.AssignRole(cmd.Role)
	}
	if cmd.IsActive != nil {
		user.IsActive = *cmd.IsActive
	}

	// Save changes
	if err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

