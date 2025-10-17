package domain

import "context"

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	// Payment operations
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, paymentID string) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int, status string) ([]*Payment, int, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	Delete(ctx context.Context, paymentID string) error
	
	// Payment status operations
	UpdateStatus(ctx context.Context, paymentID string, status PaymentStatus) error
	GetByStatus(ctx context.Context, status PaymentStatus, offset, limit int) ([]*Payment, int, error)
	
	// Payment statistics
	GetPaymentStats(ctx context.Context, userID *uint, startDate, endDate *string) (*PaymentStats, error)
	GetTotalAmountByStatus(ctx context.Context, status PaymentStatus, startDate, endDate *string) (float64, error)
	
	// Expired payments
	GetExpiredPayments(ctx context.Context) ([]*Payment, error)
	CleanupExpiredPayments(ctx context.Context) (int, error)
}

// PaymentMethodRepository defines the interface for payment method data operations
type PaymentMethodRepository interface {
	// Payment method operations
	Create(ctx context.Context, paymentMethod *PaymentMethodInfo) error
	GetByID(ctx context.Context, paymentMethodID string) (*PaymentMethodInfo, error)
	GetByUserID(ctx context.Context, userID uint) ([]*PaymentMethodInfo, error)
	GetDefaultByUserID(ctx context.Context, userID uint) (*PaymentMethodInfo, error)
	Update(ctx context.Context, paymentMethod *PaymentMethodInfo) error
	Delete(ctx context.Context, paymentMethodID string) error
	
	// Payment method status operations
	SetDefault(ctx context.Context, userID uint, paymentMethodID string) error
	SetActive(ctx context.Context, paymentMethodID string, isActive bool) error
	SetAllNonDefault(ctx context.Context, userID uint) error
}

// RefundRepository defines the interface for refund data operations
type RefundRepository interface {
	// Refund operations
	Create(ctx context.Context, refund *Refund) error
	GetByID(ctx context.Context, refundID string) (*Refund, error)
	GetByPaymentID(ctx context.Context, paymentID string) ([]*Refund, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*Refund, int, error)
	Update(ctx context.Context, refund *Refund) error
	Delete(ctx context.Context, refundID string) error
	
	// Refund status operations
	UpdateStatus(ctx context.Context, refundID string, status string) error
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*Refund, int, error)
	
	// Refund statistics
	GetRefundStats(ctx context.Context, userID *uint, startDate, endDate *string) (*RefundStats, error)
	GetTotalRefundAmount(ctx context.Context, userID *uint, startDate, endDate *string) (float64, error)
}

// PaymentStats represents payment statistics
type PaymentStats struct {
	TotalPayments      int     `json:"total_payments"`
	TotalAmount        float64 `json:"total_amount"`
	SuccessfulPayments int     `json:"successful_payments"`
	FailedPayments     int     `json:"failed_payments"`
	PendingPayments    int     `json:"pending_payments"`
	RefundedAmount     float64 `json:"refunded_amount"`
	AverageAmount      float64 `json:"average_amount"`
}

// RefundStats represents refund statistics
type RefundStats struct {
	TotalRefunds       int     `json:"total_refunds"`
	TotalAmount        float64 `json:"total_amount"`
	CompletedRefunds   int     `json:"completed_refunds"`
	PendingRefunds     int     `json:"pending_refunds"`
	FailedRefunds      int     `json:"failed_refunds"`
	AverageAmount      float64 `json:"average_amount"`
}
