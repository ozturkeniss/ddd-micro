package command

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/user/domain"
)

var (
	ErrInvalidRole = errors.New("invalid role")
)

// AssignRoleCommand represents a command to assign a role to a user
type AssignRoleCommand struct {
	UserID uint
	Role   domain.Role
}

// AssignRoleHandler handles the AssignRoleCommand
type AssignRoleHandler struct {
	repo domain.UserRepository
}

// NewAssignRoleHandler creates a new AssignRoleHandler
func NewAssignRoleHandler(repo domain.UserRepository) *AssignRoleHandler {
	return &AssignRoleHandler{
		repo: repo,
	}
}

// Handle executes the AssignRoleCommand
func (h *AssignRoleHandler) Handle(ctx context.Context, cmd AssignRoleCommand) (*domain.User, error) {
	// Validate role
	if !cmd.Role.IsValid() {
		return nil, ErrInvalidRole
	}

	// Get user
	user, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Assign role
	user.AssignRole(cmd.Role)

	// Save changes
	if err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
