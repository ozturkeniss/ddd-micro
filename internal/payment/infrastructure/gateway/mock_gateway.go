package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/payment/domain"
	"github.com/google/uuid"
)

// mockGateway implements domain.PaymentGateway for development/testing
type mockGateway struct{}

// NewMockGateway creates a new mock payment gateway
func NewMockGateway() domain.PaymentGateway {
	return &mockGateway{}
}

// CreatePayment creates a mock payment
func (g *mockGateway) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.GatewayResponse, error) {
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	transactionID := uuid.New().String()
	paymentURL := fmt.Sprintf("https://mock-payment.com/checkout/%s", transactionID)

	return &domain.GatewayResponse{
		TransactionID: transactionID,
		PaymentURL:    paymentURL,
		ClientSecret:  "mock_client_secret_" + transactionID,
		GatewayResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"payment_url":    paymentURL,
			"status":         "pending",
		},
	}, nil
}

// ProcessPayment processes a mock payment
func (g *mockGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentMethodID string) (*domain.GatewayResponse, error) {
	// Simulate some processing time
	time.Sleep(200 * time.Millisecond)

	// Simulate 90% success rate
	success := time.Now().UnixNano()%10 != 0

	var status domain.PaymentStatus
	if success {
		status = domain.PaymentStatusCompleted
	} else {
		status = domain.PaymentStatusFailed
	}

	transactionID := uuid.New().String()
	if payment.TransactionID != nil {
		transactionID = *payment.TransactionID
	}

	return &domain.GatewayResponse{
		TransactionID: transactionID,
		Status:        status,
		GatewayResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"status":         string(status),
			"success":        success,
		},
	}, nil
}

// CancelPayment cancels a mock payment
func (g *mockGateway) CancelPayment(ctx context.Context, payment *domain.Payment) (*domain.GatewayResponse, error) {
	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	transactionID := uuid.New().String()
	if payment.TransactionID != nil {
		transactionID = *payment.TransactionID
	}

	return &domain.GatewayResponse{
		TransactionID: transactionID,
		Status:        domain.PaymentStatusCancelled,
		GatewayResponse: map[string]interface{}{
			"transaction_id": transactionID,
			"status":         "cancelled",
		},
	}, nil
}

// RefundPayment refunds a mock payment
func (g *mockGateway) RefundPayment(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*domain.GatewayResponse, error) {
	// Simulate some processing time
	time.Sleep(300 * time.Millisecond)

	// Simulate 95% success rate for refunds
	success := time.Now().UnixNano()%20 != 0

	var status domain.PaymentStatus
	if success {
		status = domain.PaymentStatusRefunded
	} else {
		status = domain.PaymentStatusFailed
	}

	refundID := uuid.New().String()

	return &domain.GatewayResponse{
		TransactionID: refundID,
		Status:        status,
		GatewayResponse: map[string]interface{}{
			"refund_id": refundID,
			"status":    string(status),
			"amount":    amount,
			"reason":    reason,
			"success":   success,
		},
	}, nil
}

// CreatePaymentMethod creates a mock payment method
func (g *mockGateway) CreatePaymentMethod(ctx context.Context, userID uint, paymentMethod *domain.PaymentMethodInfo) (*domain.GatewayResponse, error) {
	// Simulate some processing time
	time.Sleep(150 * time.Millisecond)

	paymentMethodID := uuid.New().String()
	customerID := fmt.Sprintf("cus_mock_%d", userID)

	// Update payment method with mock data
	paymentMethod.Provider = "mock"
	paymentMethod.Token = paymentMethodID
	paymentMethod.LastFourDigits = stringPtr("1234")
	paymentMethod.ExpiryMonth = intPtr(12)
	paymentMethod.ExpiryYear = intPtr(2025)

	return &domain.GatewayResponse{
		TransactionID: paymentMethodID,
		GatewayResponse: map[string]interface{}{
			"payment_method_id": paymentMethodID,
			"customer_id":       customerID,
			"last_four":         "1234",
			"expiry_month":      12,
			"expiry_year":       2025,
		},
	}, nil
}

// ValidateWebhook validates a mock webhook
func (g *mockGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (map[string]interface{}, error) {
	// Simulate webhook validation
	time.Sleep(50 * time.Millisecond)

	return map[string]interface{}{
		"type": "payment_intent.succeeded",
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id":     "pi_mock_" + uuid.New().String(),
				"status": "succeeded",
			},
		},
	}, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
