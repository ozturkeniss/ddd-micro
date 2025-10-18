package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/payment/domain"
	"gorm.io/gorm"
)

// refundRepository implements domain.RefundRepository
type refundRepository struct {
	db *gorm.DB
}

// NewRefundRepository creates a new refund repository
func NewRefundRepository(db *gorm.DB) domain.RefundRepository {
	return &refundRepository{
		db: db,
	}
}

// Create creates a new refund
func (r *refundRepository) Create(ctx context.Context, refund *domain.Refund) error {
	if err := r.db.WithContext(ctx).Create(refund).Error; err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}
	return nil
}

// GetByID gets a refund by ID
func (r *refundRepository) GetByID(ctx context.Context, refundID string) (*domain.Refund, error) {
	var refund domain.Refund
	if err := r.db.WithContext(ctx).Where("id = ?", refundID).First(&refund).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRefundNotFound
		}
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}
	return &refund, nil
}

// GetByPaymentID gets refunds by payment ID
func (r *refundRepository) GetByPaymentID(ctx context.Context, paymentID string) ([]*domain.Refund, error) {
	var refunds []*domain.Refund
	if err := r.db.WithContext(ctx).Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&refunds).Error; err != nil {
		return nil, fmt.Errorf("failed to get refunds by payment ID: %w", err)
	}
	return refunds, nil
}

// GetByUserID gets refunds by user ID with pagination
func (r *refundRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*domain.Refund, int, error) {
	var refunds []*domain.Refund
	var total int64

	// Join with payments table to filter by user_id
	query := r.db.WithContext(ctx).Model(&domain.Refund{}).
		Joins("JOIN payments ON refunds.payment_id = payments.id").
		Where("payments.user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count refunds: %w", err)
	}

	// Get refunds with pagination
	if err := query.Order("refunds.created_at DESC").Limit(limit).Offset(offset).Find(&refunds).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get refunds: %w", err)
	}

	return refunds, int(total), nil
}

// Update updates a refund
func (r *refundRepository) Update(ctx context.Context, refund *domain.Refund) error {
	refund.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Save(refund).Error; err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}
	return nil
}

// Delete deletes a refund
func (r *refundRepository) Delete(ctx context.Context, refundID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", refundID).Delete(&domain.Refund{}).Error; err != nil {
		return fmt.Errorf("failed to delete refund: %w", err)
	}
	return nil
}

// UpdateStatus updates refund status
func (r *refundRepository) UpdateStatus(ctx context.Context, refundID string, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "completed" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := r.db.WithContext(ctx).Model(&domain.Refund{}).Where("id = ?", refundID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update refund status: %w", err)
	}
	return nil
}

// GetByStatus gets refunds by status with pagination
func (r *refundRepository) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.Refund, int, error) {
	var refunds []*domain.Refund
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Refund{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count refunds by status: %w", err)
	}

	// Get refunds with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&refunds).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get refunds by status: %w", err)
	}

	return refunds, int(total), nil
}

// GetRefundStats gets refund statistics
func (r *refundRepository) GetRefundStats(ctx context.Context, userID *uint, startDate, endDate *string) (*domain.RefundStats, error) {
	var stats domain.RefundStats

	query := r.db.WithContext(ctx).Model(&domain.Refund{})

	// Apply user filter if provided
	if userID != nil {
		query = query.Joins("JOIN payments ON refunds.payment_id = payments.id").
			Where("payments.user_id = ?", *userID)
	}

	// Apply date range filter if provided
	if startDate != nil && endDate != nil {
		query = query.Where("refunds.created_at BETWEEN ? AND ?", *startDate, *endDate)
	}

	// Get total refunds and amount
	if err := query.Count(&stats.TotalRefunds).Error; err != nil {
		return nil, fmt.Errorf("failed to count total refunds: %w", err)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total refund amount: %w", err)
	}

	// Get completed refunds
	if err := query.Where("status = ?", "completed").Count(&stats.CompletedRefunds).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed refunds: %w", err)
	}

	// Get pending refunds
	if err := query.Where("status = ?", "pending").Count(&stats.PendingRefunds).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending refunds: %w", err)
	}

	// Get failed refunds
	if err := query.Where("status = ?", "failed").Count(&stats.FailedRefunds).Error; err != nil {
		return nil, fmt.Errorf("failed to count failed refunds: %w", err)
	}

	// Calculate average amount
	if stats.TotalRefunds > 0 {
		stats.AverageAmount = stats.TotalAmount / float64(stats.TotalRefunds)
	}

	return &stats, nil
}

// GetTotalRefundAmount gets total refund amount
func (r *refundRepository) GetTotalRefundAmount(ctx context.Context, userID *uint, startDate, endDate *string) (float64, error) {
	var total float64

	query := r.db.WithContext(ctx).Model(&domain.Refund{})

	// Apply user filter if provided
	if userID != nil {
		query = query.Joins("JOIN payments ON refunds.payment_id = payments.id").
			Where("payments.user_id = ?", *userID)
	}

	// Apply date range filter if provided
	if startDate != nil && endDate != nil {
		query = query.Where("refunds.created_at BETWEEN ? AND ?", *startDate, *endDate)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to get total refund amount: %w", err)
	}

	return total, nil
}
