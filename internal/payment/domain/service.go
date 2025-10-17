package domain

import "context"

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	// Payment operations
	CreatePayment(ctx context.Context, userID uint, orderID string, amount float64, currency string, paymentMethod string) (*Payment, error)
	ProcessPayment(ctx context.Context, paymentID string, paymentMethodID string) (*Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	GetUserPayments(ctx context.Context, userID uint, offset, limit int) ([]*Payment, int, error)
	CancelPayment(ctx context.Context, paymentID string) (*Payment, error)
	
	// Payment method operations
	AddPaymentMethod(ctx context.Context, userID uint, paymentMethod *PaymentMethodInfo) (*PaymentMethodInfo, error)
	GetPaymentMethods(ctx context.Context, userID uint) ([]*PaymentMethodInfo, error)
	GetDefaultPaymentMethod(ctx context.Context, userID uint) (*PaymentMethodInfo, error)
	UpdatePaymentMethod(ctx context.Context, paymentMethodID string, paymentMethod *PaymentMethodInfo) (*PaymentMethodInfo, error)
	DeletePaymentMethod(ctx context.Context, paymentMethodID string) error
	SetDefaultPaymentMethod(ctx context.Context, userID uint, paymentMethodID string) error
	
	// Refund operations
	CreateRefund(ctx context.Context, paymentID string, amount float64, reason string) (*Refund, error)
	GetRefund(ctx context.Context, refundID string) (*Refund, error)
	GetRefundsByPayment(ctx context.Context, paymentID string) ([]*Refund, error)
	GetUserRefunds(ctx context.Context, userID uint, offset, limit int) ([]*Refund, int, error)
	ProcessRefund(ctx context.Context, refundID string) (*Refund, error)
	
	// Webhook operations
	ProcessWebhook(ctx context.Context, provider string, payload []byte, signature string) error
	
	// Statistics
	GetPaymentStats(ctx context.Context, userID *uint, startDate, endDate *string) (*PaymentStats, error)
	GetRefundStats(ctx context.Context, userID *uint, startDate, endDate *string) (*RefundStats, error)
	
	// Admin operations
	GetAllPayments(ctx context.Context, offset, limit int, filters *PaymentFilters) ([]*Payment, int, error)
	GetAllRefunds(ctx context.Context, offset, limit int, filters *RefundFilters) ([]*Refund, int, error)
	CleanupExpiredPayments(ctx context.Context) (int, error)
}

// PaymentFilters represents filters for payment queries
type PaymentFilters struct {
	UserID        *uint           `json:"user_id,omitempty"`
	Status        *PaymentStatus  `json:"status,omitempty"`
	PaymentMethod *PaymentMethod  `json:"payment_method,omitempty"`
	StartDate     *string         `json:"start_date,omitempty"`
	EndDate       *string         `json:"end_date,omitempty"`
	MinAmount     *float64        `json:"min_amount,omitempty"`
	MaxAmount     *float64        `json:"max_amount,omitempty"`
}

// RefundFilters represents filters for refund queries
type RefundFilters struct {
	UserID    *uint   `json:"user_id,omitempty"`
	Status    *string `json:"status,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	MinAmount *float64 `json:"min_amount,omitempty"`
	MaxAmount *float64 `json:"max_amount,omitempty"`
}

// PaymentNotificationService defines the interface for payment notifications
type PaymentNotificationService interface {
	// Send payment notifications
	SendPaymentConfirmation(ctx context.Context, payment *Payment) error
	SendPaymentFailed(ctx context.Context, payment *Payment) error
	SendRefundConfirmation(ctx context.Context, refund *Refund) error
	
	// Send admin notifications
	SendPaymentAlert(ctx context.Context, payment *Payment, alertType string) error
	SendRefundAlert(ctx context.Context, refund *Refund, alertType string) error
}

// PaymentAuditService defines the interface for payment auditing
type PaymentAuditService interface {
	// Log payment events
	LogPaymentCreated(ctx context.Context, payment *Payment) error
	LogPaymentUpdated(ctx context.Context, payment *Payment, oldStatus PaymentStatus) error
	LogPaymentCompleted(ctx context.Context, payment *Payment) error
	LogPaymentFailed(ctx context.Context, payment *Payment, reason string) error
	LogRefundCreated(ctx context.Context, refund *Refund) error
	LogRefundProcessed(ctx context.Context, refund *Refund) error
	
	// Get audit logs
	GetPaymentAuditLogs(ctx context.Context, paymentID string) ([]*AuditLog, error)
	GetUserAuditLogs(ctx context.Context, userID uint, offset, limit int) ([]*AuditLog, int, error)
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID          string                 `json:"id"`
	EntityType  string                 `json:"entity_type"` // payment, refund, payment_method
	EntityID    string                 `json:"entity_id"`
	UserID      uint                   `json:"user_id"`
	Action      string                 `json:"action"`
	OldValues   map[string]interface{} `json:"old_values,omitempty"`
	NewValues   map[string]interface{} `json:"new_values,omitempty"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	CreatedAt   string                 `json:"created_at"`
}
