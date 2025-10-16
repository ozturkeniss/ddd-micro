package domain

import "context"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uint) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// Exists checks if a user exists by email
	Exists(ctx context.Context, email string) (bool, error)
}
