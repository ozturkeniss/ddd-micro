package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/payment/domain"
	"gorm.io/gorm"
)

// paymentMethodRepository implements domain.PaymentMethodRepository
type paymentMethodRepository struct {
	db *gorm.DB
}

// NewPaymentMethodRepository creates a new payment method repository
func NewPaymentMethodRepository(db *gorm.DB) domain.PaymentMethodRepository {
	return &paymentMethodRepository{
		db: db,
	}
}

// Create creates a new payment method
func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *domain.PaymentMethodInfo) error {
	if err := r.db.WithContext(ctx).Create(paymentMethod).Error; err != nil {
		return fmt.Errorf("failed to create payment method: %w", err)
	}
	return nil
}

// GetByID gets a payment method by ID
func (r *paymentMethodRepository) GetByID(ctx context.Context, paymentMethodID string) (*domain.PaymentMethodInfo, error) {
	var paymentMethod domain.PaymentMethodInfo
	if err := r.db.WithContext(ctx).Where("id = ?", paymentMethodID).First(&paymentMethod).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentMethodNotFound
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	return &paymentMethod, nil
}

// GetByUserID gets payment methods by user ID
func (r *paymentMethodRepository) GetByUserID(ctx context.Context, userID uint) ([]*domain.PaymentMethodInfo, error) {
	var paymentMethods []*domain.PaymentMethodInfo
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Order("is_default DESC, created_at DESC").Find(&paymentMethods).Error; err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	return paymentMethods, nil
}

// GetDefaultByUserID gets the default payment method for a user
func (r *paymentMethodRepository) GetDefaultByUserID(ctx context.Context, userID uint) (*domain.PaymentMethodInfo, error) {
	var paymentMethod domain.PaymentMethodInfo
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_default = ? AND is_active = ?", userID, true, true).First(&paymentMethod).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentMethodNotFound
		}
		return nil, fmt.Errorf("failed to get default payment method: %w", err)
	}
	return &paymentMethod, nil
}

// Update updates a payment method
func (r *paymentMethodRepository) Update(ctx context.Context, paymentMethod *domain.PaymentMethodInfo) error {
	paymentMethod.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(paymentMethod).Error; err != nil {
		return fmt.Errorf("failed to update payment method: %w", err)
	}
	return nil
}

// Delete deletes a payment method
func (r *paymentMethodRepository) Delete(ctx context.Context, paymentMethodID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", paymentMethodID).Delete(&domain.PaymentMethodInfo{}).Error; err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}
	return nil
}

// SetDefault sets a payment method as default for a user
func (r *paymentMethodRepository) SetDefault(ctx context.Context, userID uint, paymentMethodID string) error {
	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Set all payment methods for this user to not be default
	if err := tx.Model(&domain.PaymentMethodInfo{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set all payment methods as non-default: %w", err)
	}

	// Set the specified payment method as default
	if err := tx.Model(&domain.PaymentMethodInfo{}).Where("id = ? AND user_id = ?", paymentMethodID, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set payment method as default: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SetActive sets payment method active status
func (r *paymentMethodRepository) SetActive(ctx context.Context, paymentMethodID string, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&domain.PaymentMethodInfo{}).Where("id = ?", paymentMethodID).Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to set payment method active status: %w", err)
	}
	return nil
}

// SetAllNonDefault sets all payment methods for a user as non-default
func (r *paymentMethodRepository) SetAllNonDefault(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&domain.PaymentMethodInfo{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		return fmt.Errorf("failed to set all payment methods as non-default: %w", err)
	}
	return nil
}
