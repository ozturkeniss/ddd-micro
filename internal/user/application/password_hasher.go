package application

import "github.com/ddd-micro/pkg/security"

// PasswordHasher is an alias for backward compatibility
type PasswordHasher = security.PasswordHasher

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher() *PasswordHasher {
	return security.NewPasswordHasher()
}

