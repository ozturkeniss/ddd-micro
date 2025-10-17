package dto

import "time"

// Payment DTOs

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	OrderID         string  `json:"order_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,min=0.01"`
	Currency        string  `json:"currency" binding:"required,len=3"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	PaymentMethodID string  `json:"payment_method_id,omitempty"`
	ReturnURL       string  `json:"return_url,omitempty"`
	CancelURL       string  `json:"cancel_url,omitempty"`
}

// ProcessPaymentRequest represents the request to process a payment
type ProcessPaymentRequest struct {
	PaymentMethodID string                 `json:"payment_method_id" binding:"required"`
	ConfirmationData map[string]interface{} `json:"confirmation_data,omitempty"`
}

// PaymentResponse represents the response for payment operations
type PaymentResponse struct {
	ID              string                 `json:"id"`
	UserID          uint                   `json:"user_id"`
	OrderID         string                 `json:"order_id"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	Status          string                 `json:"status"`
	PaymentMethod   string                 `json:"payment_method"`
	PaymentProvider string                 `json:"payment_provider"`
	TransactionID   *string                `json:"transaction_id,omitempty"`
	GatewayResponse map[string]interface{} `json:"gateway_response,omitempty"`
	PaymentURL      string                 `json:"payment_url,omitempty"`
	ClientSecret    string                 `json:"client_secret,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
}

// ListPaymentsRequest represents the request for listing payments
type ListPaymentsRequest struct {
	UserID uint   `json:"user_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status,omitempty"`
}

// ListPaymentsResponse represents the response for listing payments
type PaymentListResponse struct {
	Payments   []PaymentResponse `json:"payments"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
	HasNext    bool              `json:"has_next"`
	HasPrev    bool              `json:"has_prev"`
}

// PaymentMethod DTOs

// AddPaymentMethodRequest represents the request to add a payment method
type AddPaymentMethodRequest struct {
	Type            string `json:"type" binding:"required"`
	Provider        string `json:"provider" binding:"required"`
	Token           string `json:"token" binding:"required"`
	IsDefault       bool   `json:"is_default"`
}

// UpdatePaymentMethodRequest represents the request to update a payment method
type UpdatePaymentMethodRequest struct {
	IsDefault bool `json:"is_default"`
	IsActive  bool `json:"is_active"`
}

// PaymentMethodResponse represents the response for payment method operations
type PaymentMethodResponse struct {
	ID             string    `json:"id"`
	UserID         uint      `json:"user_id"`
	Type           string    `json:"type"`
	Provider       string    `json:"provider"`
	LastFourDigits *string   `json:"last_four_digits,omitempty"`
	ExpiryMonth    *int      `json:"expiry_month,omitempty"`
	ExpiryYear     *int      `json:"expiry_year,omitempty"`
	IsDefault      bool      `json:"is_default"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PaymentMethodListResponse represents the response for listing payment methods
type PaymentMethodListResponse struct {
	PaymentMethods []PaymentMethodResponse `json:"payment_methods"`
	Total          int                     `json:"total"`
}

// Refund DTOs

// CreateRefundRequest represents the request to create a refund
type CreateRefundRequest struct {
	PaymentID string  `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0.01"`
	Reason    string  `json:"reason" binding:"required"`
}

// RefundResponse represents the response for refund operations
type RefundResponse struct {
	ID          string     `json:"id"`
	PaymentID   string     `json:"payment_id"`
	Amount      float64    `json:"amount"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// AdminListPaymentsRequest represents the request for admin listing payments
type AdminListPaymentsRequest struct {
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
	UserID *uint   `json:"user_id,omitempty"`
	Status string  `json:"status,omitempty"`
}

// UpdatePaymentStatusRequest represents the request to update payment status
type UpdatePaymentStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason,omitempty"`
}

// AdminListRefundsRequest represents the request for admin listing refunds
type AdminListRefundsRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	PaymentID string `json:"payment_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

// RefundListResponse represents the response for listing refunds
type RefundListResponse struct {
	Refunds []RefundResponse `json:"refunds"`
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}

// Statistics DTOs

// PaymentStatsResponse represents the response for payment statistics
type PaymentStatsResponse struct {
	TotalPayments      int     `json:"total_payments"`
	TotalAmount        float64 `json:"total_amount"`
	SuccessfulPayments int     `json:"successful_payments"`
	FailedPayments     int     `json:"failed_payments"`
	PendingPayments    int     `json:"pending_payments"`
	RefundedAmount     float64 `json:"refunded_amount"`
	AverageAmount      float64 `json:"average_amount"`
}

// RefundStatsResponse represents the response for refund statistics
type RefundStatsResponse struct {
	TotalRefunds     int     `json:"total_refunds"`
	TotalAmount      float64 `json:"total_amount"`
	CompletedRefunds int     `json:"completed_refunds"`
	PendingRefunds   int     `json:"pending_refunds"`
	FailedRefunds    int     `json:"failed_refunds"`
	AverageAmount    float64 `json:"average_amount"`
}

// Webhook DTOs

// WebhookRequest represents the request for webhook processing
type WebhookRequest struct {
	Provider  string `json:"provider" binding:"required"`
	Payload   string `json:"payload" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// WebhookResponse represents the response for webhook processing
type WebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Common DTOs

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationRequest represents a pagination request
type PaginationRequest struct {
	Offset int `json:"offset" form:"offset" binding:"min=0"`
	Limit  int `json:"limit" form:"limit" binding:"min=1,max=100"`
}

// DateRangeRequest represents a date range request
type DateRangeRequest struct {
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
}
