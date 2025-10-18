package domain

import (
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

// PaymentMethod represents the payment method type
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodPayPal     PaymentMethod = "paypal"
	PaymentMethodStripe     PaymentMethod = "stripe"
)

// Payment represents a payment transaction
type Payment struct {
	ID               string                 `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID           uint                   `json:"user_id" gorm:"not null;index"`
	OrderID          string                 `json:"order_id" gorm:"not null;index;type:varchar(36)"`
	Amount           float64                `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency         string                 `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	Status           PaymentStatus          `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentMethod    PaymentMethod          `json:"payment_method" gorm:"type:varchar(20);not null"`
	PaymentProvider  string                 `json:"payment_provider" gorm:"type:varchar(50);not null"`
	TransactionID    *string                `json:"transaction_id" gorm:"type:varchar(100);index"`
	GatewayResponse  map[string]interface{} `json:"gateway_response" gorm:"type:jsonb"`
	ReturnURL        *string                `json:"return_url" gorm:"type:text"`
	CancelURL        *string                `json:"cancel_url" gorm:"type:text"`
	// Optional: Direct product purchase (without basket)
	ProductID        *uint                  `json:"product_id" gorm:"index"`
	Quantity         *int                   `json:"quantity"`
	// Optional: Basket-based purchase
	BasketID         *string                `json:"basket_id" gorm:"type:varchar(36);index"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	CompletedAt      *time.Time             `json:"completed_at"`
	ExpiresAt        *time.Time             `json:"expires_at" gorm:"index"`
}

// PaymentMethodInfo represents a user's payment method
type PaymentMethodInfo struct {
	ID              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID          uint      `json:"user_id" gorm:"not null;index"`
	Type            string    `json:"type" gorm:"type:varchar(20);not null"`
	Provider        string    `json:"provider" gorm:"type:varchar(50);not null"`
	LastFourDigits  *string   `json:"last_four_digits" gorm:"type:varchar(4)"`
	ExpiryMonth     *int      `json:"expiry_month" gorm:"type:smallint"`
	ExpiryYear      *int      `json:"expiry_year" gorm:"type:smallint"`
	IsDefault       bool      `json:"is_default" gorm:"default:false"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Refund represents a refund transaction
type Refund struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	PaymentID   string    `json:"payment_id" gorm:"not null;index;type:varchar(36)"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Reason      string    `json:"reason" gorm:"type:text;not null"`
	Status      string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// TableName returns the table name for Payment
func (Payment) TableName() string {
	return "payments"
}

// TableName returns the table name for PaymentMethodInfo
func (PaymentMethodInfo) TableName() string {
	return "payment_methods"
}

// TableName returns the table name for Refund
func (Refund) TableName() string {
	return "refunds"
}

// IsCompleted checks if the payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsFailed checks if the payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// IsPending checks if the payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending
}

// IsProcessing checks if the payment is processing
func (p *Payment) IsProcessing() bool {
	return p.Status == PaymentStatusProcessing
}

// IsCancelled checks if the payment is cancelled
func (p *Payment) IsCancelled() bool {
	return p.Status == PaymentStatusCancelled
}

// IsRefunded checks if the payment is refunded
func (p *Payment) IsRefunded() bool {
	return p.Status == PaymentStatusRefunded
}

// CanBeRefunded checks if the payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.IsCompleted() && !p.IsRefunded()
}

// CanBeCancelled checks if the payment can be cancelled
func (p *Payment) CanBeCancelled() bool {
	return p.IsPending() || p.IsProcessing()
}

// SetCompleted marks the payment as completed
func (p *Payment) SetCompleted() {
	p.Status = PaymentStatusCompleted
	now := time.Now()
	p.CompletedAt = &now
}

// SetFailed marks the payment as failed
func (p *Payment) SetFailed() {
	p.Status = PaymentStatusFailed
}

// SetProcessing marks the payment as processing
func (p *Payment) SetProcessing() {
	p.Status = PaymentStatusProcessing
}

// SetCancelled marks the payment as cancelled
func (p *Payment) SetCancelled() {
	p.Status = PaymentStatusCancelled
}

// SetRefunded marks the payment as refunded
func (p *Payment) SetRefunded() {
	p.Status = PaymentStatusRefunded
}

// SetExpiration sets the expiration time for the payment
func (p *Payment) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().Add(duration)
	p.ExpiresAt = &expiresAt
}

// IsExpired checks if the payment has expired
func (p *Payment) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

// Validate validates the payment
func (p *Payment) Validate() error {
	if p.UserID == 0 {
		return ErrInvalidUserID
	}
	
	if p.OrderID == "" {
		return ErrInvalidOrderID
	}
	
	if p.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	if p.Currency == "" {
		return ErrInvalidCurrency
	}
	
	if p.PaymentMethod == "" {
		return ErrInvalidPaymentMethod
	}
	
	if p.PaymentProvider == "" {
		return ErrInvalidPaymentProvider
	}
	
	return nil
}

// Validate validates the payment method
func (pm *PaymentMethodInfo) Validate() error {
	if pm.UserID == 0 {
		return ErrInvalidUserID
	}
	
	if pm.Type == "" {
		return ErrInvalidPaymentMethodType
	}
	
	if pm.Provider == "" {
		return ErrInvalidPaymentProvider
	}
	
	return nil
}

// Validate validates the refund
func (r *Refund) Validate() error {
	if r.PaymentID == "" {
		return ErrInvalidPaymentID
	}
	
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	if r.Reason == "" {
		return ErrInvalidRefundReason
	}
	
	return nil
}
