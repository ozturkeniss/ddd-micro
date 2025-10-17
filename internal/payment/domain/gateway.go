package domain

import "context"

// PaymentGateway defines the interface for payment gateway operations
type PaymentGateway interface {
	// Payment operations
	CreatePayment(ctx context.Context, payment *Payment) (*PaymentGatewayResponse, error)
	ProcessPayment(ctx context.Context, payment *Payment, paymentMethodID string) (*PaymentGatewayResponse, error)
	CancelPayment(ctx context.Context, payment *Payment) (*PaymentGatewayResponse, error)
	RefundPayment(ctx context.Context, payment *Payment, amount float64, reason string) (*PaymentGatewayResponse, error)
	
	// Payment method operations
	CreatePaymentMethod(ctx context.Context, userID uint, paymentMethod *PaymentMethodInfo) (*PaymentGatewayResponse, error)
	DeletePaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentGatewayResponse, error)
	
	// Webhook operations
	ProcessWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)
	
	// Health check
	HealthCheck(ctx context.Context) error
}

// PaymentGatewayResponse represents the response from payment gateway
type PaymentGatewayResponse struct {
	Success         bool                   `json:"success"`
	TransactionID   string                 `json:"transaction_id"`
	PaymentURL      string                 `json:"payment_url"`
	ClientSecret    string                 `json:"client_secret"`
	Status          PaymentStatus          `json:"status"`
	Message         string                 `json:"message"`
	GatewayResponse map[string]interface{} `json:"gateway_response"`
	Error           *PaymentGatewayError   `json:"error,omitempty"`
}

// PaymentGatewayError represents an error from payment gateway
type PaymentGatewayError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Param       string `json:"param,omitempty"`
	DeclineCode string `json:"decline_code,omitempty"`
}

// WebhookEvent represents a webhook event from payment gateway
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Created   int64                  `json:"created"`
	Processed bool                   `json:"processed"`
}

// PaymentGatewayConfig represents configuration for payment gateway
type PaymentGatewayConfig struct {
	Provider     string `json:"provider"`
	APIKey       string `json:"api_key"`
	SecretKey    string `json:"secret_key"`
	WebhookSecret string `json:"webhook_secret"`
	Environment  string `json:"environment"` // sandbox, production
	BaseURL      string `json:"base_url"`
	Timeout      int    `json:"timeout"` // in seconds
}

// PaymentGatewayFactory creates payment gateway instances
type PaymentGatewayFactory interface {
	CreateGateway(config *PaymentGatewayConfig) (PaymentGateway, error)
	GetSupportedProviders() []string
}

// Supported payment gateway providers
const (
	ProviderStripe = "stripe"
	ProviderPayPal = "paypal"
	ProviderSquare = "square"
	ProviderRazorpay = "razorpay"
)

// Payment gateway event types
const (
	EventPaymentSucceeded = "payment.succeeded"
	EventPaymentFailed    = "payment.failed"
	EventPaymentCancelled = "payment.cancelled"
	EventRefundSucceeded  = "refund.succeeded"
	EventRefundFailed     = "refund.failed"
)
