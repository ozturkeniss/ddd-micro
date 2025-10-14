package persistence

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/user/domain"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

// UserRepository is the concrete implementation of domain.UserRepository
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	// Check if user already exists
	exists, err := r.Exists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExists
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).First(&user, id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, id)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	
	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&users)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// Exists checks if a user exists by email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("email = ?", email).
		Count(&count)
	
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

