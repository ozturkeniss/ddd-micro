package command

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
)

var (
	ErrInvalidOldPassword = errors.New("invalid old password")
)

// ChangePasswordCommand represents a command to change user password
type ChangePasswordCommand struct {
	UserID      uint
	OldPassword string
	NewPassword string
}

// ChangePasswordHandler handles the ChangePasswordCommand
type ChangePasswordHandler struct {
	repo           domain.UserRepository
	passwordHasher *application.PasswordHasher
}

// NewChangePasswordHandler creates a new ChangePasswordHandler
func NewChangePasswordHandler(repo domain.UserRepository) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		repo:           repo,
		passwordHasher: application.NewPasswordHasher(),
	}
}

// Handle executes the ChangePasswordCommand
func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	// Get user
	user, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := h.passwordHasher.ComparePassword(user.Password, cmd.OldPassword); err != nil {
		return ErrInvalidOldPassword
	}

	// Hash new password
	hashedPassword, err := h.passwordHasher.HashPassword(cmd.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	return h.repo.Update(ctx, user)
}

