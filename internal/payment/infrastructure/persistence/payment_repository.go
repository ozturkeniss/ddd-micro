package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/payment/domain"
	"gorm.io/gorm"
)

// paymentRepository implements domain.PaymentRepository
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

// Create creates a new payment
func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

// GetByID gets a payment by ID
func (r *paymentRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", paymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &payment, nil
}

// GetByOrderID gets a payment by order ID
func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by order ID: %w", err)
	}
	return &payment, nil
}

// GetByUserID gets payments by user ID with pagination and status filter
func (r *paymentRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int, status string) ([]*domain.Payment, int, error) {
	var payments []*domain.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Payment{}).Where("user_id = ?", userID)

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	// Get payments with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get payments: %w", err)
	}

	return payments, int(total), nil
}

// GetByTransactionID gets a payment by transaction ID
func (r *paymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment by transaction ID: %w", err)
	}
	return &payment, nil
}

// Update updates a payment
func (r *paymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	payment.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(payment).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return nil
}

// Delete deletes a payment
func (r *paymentRepository) Delete(ctx context.Context, paymentID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", paymentID).Delete(&domain.Payment{}).Error; err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	return nil
}

// UpdateStatus updates payment status
func (r *paymentRepository) UpdateStatus(ctx context.Context, paymentID string, status domain.PaymentStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == domain.PaymentStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := r.db.WithContext(ctx).Model(&domain.Payment{}).Where("id = ?", paymentID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}

// GetByStatus gets payments by status with pagination
func (r *paymentRepository) GetByStatus(ctx context.Context, status domain.PaymentStatus, limit, offset int) ([]*domain.Payment, int, error) {
	var payments []*domain.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Payment{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count payments by status: %w", err)
	}

	// Get payments with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get payments by status: %w", err)
	}

	return payments, int(total), nil
}

// GetPaymentStats gets payment statistics
func (r *paymentRepository) GetPaymentStats(ctx context.Context, userID *uint, startDate, endDate *string) (*domain.PaymentStats, error) {
	var stats domain.PaymentStats

	query := r.db.WithContext(ctx).Model(&domain.Payment{})

	// Apply user filter if provided
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Apply date range filter if provided
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	}

	// Get total payments and amount
	if err := query.Count(&stats.TotalPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to count total payments: %w", err)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total amount: %w", err)
	}

	// Get successful payments
	if err := query.Where("status = ?", domain.PaymentStatusCompleted).Count(&stats.SuccessfulPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to count successful payments: %w", err)
	}

	// Get failed payments
	if err := query.Where("status = ?", domain.PaymentStatusFailed).Count(&stats.FailedPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to count failed payments: %w", err)
	}

	// Get pending payments
	if err := query.Where("status = ?", domain.PaymentStatusPending).Count(&stats.PendingPayments).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending payments: %w", err)
	}

	// Get refunded amount
	if err := query.Where("status = ?", domain.PaymentStatusRefunded).Select("COALESCE(SUM(amount), 0)").Scan(&stats.RefundedAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to get refunded amount: %w", err)
	}

	// Calculate average amount
	if stats.TotalPayments > 0 {
		stats.AverageAmount = stats.TotalAmount / float64(stats.TotalPayments)
	}

	return &stats, nil
}

// GetTotalAmountByStatus gets total amount by status
func (r *paymentRepository) GetTotalAmountByStatus(ctx context.Context, status domain.PaymentStatus, startDate, endDate *string) (float64, error) {
	var total float64

	query := r.db.WithContext(ctx).Model(&domain.Payment{}).Where("status = ?", status)

	// Apply date range filter if provided
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total amount by status: %w", err)
	}

	return total, nil
}

// GetExpiredPayments gets expired payments
func (r *paymentRepository) GetExpiredPayments(ctx context.Context) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	now := time.Now()

	if err := r.db.WithContext(ctx).Where("status = ? AND expires_at < ?", domain.PaymentStatusPending, now).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get expired payments: %w", err)
	}

	return payments, nil
}

// CleanupExpiredPayments cleans up expired payments
func (r *paymentRepository) CleanupExpiredPayments(ctx context.Context) (int, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Where("status = ? AND expires_at < ?", domain.PaymentStatusPending, now).Delete(&domain.Payment{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup expired payments: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}
