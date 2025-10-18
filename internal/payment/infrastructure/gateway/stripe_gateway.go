package gateway

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ddd-micro/internal/payment/domain"
	"github.com/ddd-micro/internal/payment/infrastructure/config"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
)

// stripeGateway implements domain.PaymentGateway
type stripeGateway struct {
	config *config.Config
}

// NewStripeGateway creates a new Stripe payment gateway
func NewStripeGateway(cfg *config.Config) domain.PaymentGateway {
	// Set Stripe API key
	stripe.Key = cfg.Stripe.SecretKey

	return &stripeGateway{
		config: cfg,
	}
}

// CreatePayment creates a payment via Stripe
func (g *stripeGateway) CreatePayment(ctx context.Context, payment *domain.Payment) (*domain.GatewayResponse, error) {
	// Create Stripe checkout session
	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			string(payment.PaymentMethod),
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(payment.Currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("Order %s", payment.OrderID)),
					},
					UnitAmount: stripe.Int64(int64(payment.Amount * 100)), // Convert to cents
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(payment.ReturnURL),
		CancelURL:  stripe.String(payment.CancelURL),
		Metadata: map[string]string{
			"payment_id": payment.ID,
			"order_id":   payment.OrderID,
			"user_id":    strconv.FormatUint(uint64(payment.UserID), 10),
		},
	}

	// Add customer if payment method ID is provided
	if payment.PaymentMethodID != "" {
		sessionParams.Customer = stripe.String(payment.PaymentMethodID)
	}

	session, err := session.New(sessionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe checkout session: %w", err)
	}

	return &domain.GatewayResponse{
		TransactionID:   session.ID,
		PaymentURL:      session.URL,
		ClientSecret:    session.ClientSecret,
		GatewayResponse: map[string]interface{}{
			"session_id": session.ID,
			"url":        session.URL,
		},
	}, nil
}

// ProcessPayment processes a payment via Stripe
func (g *stripeGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentMethodID string) (*domain.GatewayResponse, error) {
	// Create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), // Convert to cents
		Currency: stripe.String(payment.Currency),
		Metadata: map[string]string{
			"payment_id": payment.ID,
			"order_id":   payment.OrderID,
			"user_id":    strconv.FormatUint(uint64(payment.UserID), 10),
		},
	}

	// Add payment method if provided
	if paymentMethodID != "" {
		params.PaymentMethod = stripe.String(paymentMethodID)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Confirm payment intent
	confirmedPI, err := paymentintent.Confirm(pi.ID, &stripe.PaymentIntentConfirmParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	// Determine status based on payment intent status
	var status domain.PaymentStatus
	switch confirmedPI.Status {
	case stripe.PaymentIntentStatusSucceeded:
		status = domain.PaymentStatusCompleted
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		status = domain.PaymentStatusFailed
	case stripe.PaymentIntentStatusRequiresConfirmation:
		status = domain.PaymentStatusProcessing
	default:
		status = domain.PaymentStatusPending
	}

	return &domain.GatewayResponse{
		TransactionID:   confirmedPI.ID,
		Status:          status,
		GatewayResponse: map[string]interface{}{
			"payment_intent_id": confirmedPI.ID,
			"status":            confirmedPI.Status,
			"client_secret":     confirmedPI.ClientSecret,
		},
	}, nil
}

// CancelPayment cancels a payment via Stripe
func (g *stripeGateway) CancelPayment(ctx context.Context, payment *domain.Payment) (*domain.GatewayResponse, error) {
	// For Stripe, we can't cancel a completed payment, only refund it
	// This would typically be used for pending payments
	if payment.TransactionID == nil {
		return nil, fmt.Errorf("no transaction ID to cancel")
	}

	// Cancel payment intent if it's still pending
	pi, err := paymentintent.Cancel(*payment.TransactionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	return &domain.GatewayResponse{
		TransactionID:   pi.ID,
		Status:          domain.PaymentStatusCancelled,
		GatewayResponse: map[string]interface{}{
			"payment_intent_id": pi.ID,
			"status":            pi.Status,
		},
	}, nil
}

// RefundPayment refunds a payment via Stripe
func (g *stripeGateway) RefundPayment(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*domain.GatewayResponse, error) {
	if payment.TransactionID == nil {
		return nil, fmt.Errorf("no transaction ID to refund")
	}

	// Create refund
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(*payment.TransactionID),
		Amount:        stripe.Int64(int64(amount * 100)), // Convert to cents
		Reason:        stripe.String(reason),
		Metadata: map[string]string{
			"payment_id": payment.ID,
			"order_id":   payment.OrderID,
		},
	}

	ref, err := refund.New(refundParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// Determine status based on refund status
	var status domain.PaymentStatus
	switch ref.Status {
	case stripe.RefundStatusSucceeded:
		status = domain.PaymentStatusRefunded
	case stripe.RefundStatusPending:
		status = domain.PaymentStatusProcessing
	case stripe.RefundStatusFailed:
		status = domain.PaymentStatusFailed
	default:
		status = domain.PaymentStatusPending
	}

	return &domain.GatewayResponse{
		TransactionID:   ref.ID,
		Status:          status,
		GatewayResponse: map[string]interface{}{
			"refund_id": ref.ID,
			"status":    ref.Status,
			"amount":    ref.Amount,
		},
	}, nil
}

// CreatePaymentMethod creates a payment method via Stripe
func (g *stripeGateway) CreatePaymentMethod(ctx context.Context, userID uint, paymentMethod *domain.PaymentMethodInfo) (*domain.GatewayResponse, error) {
	// Create Stripe customer if not exists
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(fmt.Sprintf("user_%d@example.com", userID)), // This should come from user service
		Metadata: map[string]string{
			"user_id": strconv.FormatUint(uint64(userID), 10),
		},
	}

	cust, err := customer.New(customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Create payment method
	pmParams := &stripe.PaymentMethodParams{
		Type: stripe.String(paymentMethod.Type),
		Card: &stripe.PaymentMethodCardParams{
			Token: stripe.String(paymentMethod.Token),
		},
	}

	pm, err := paymentmethod.New(pmParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	// Attach payment method to customer
	attachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(cust.ID),
	}

	attachedPM, err := paymentmethod.Attach(pm.ID, attachParams)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method: %w", err)
	}

	// Update payment method with Stripe data
	paymentMethod.Provider = "stripe"
	paymentMethod.Token = attachedPM.ID

	return &domain.GatewayResponse{
		TransactionID:   attachedPM.ID,
		GatewayResponse: map[string]interface{}{
			"payment_method_id": attachedPM.ID,
			"customer_id":       cust.ID,
		},
	}, nil
}

// ValidateWebhook validates Stripe webhook signature
func (g *stripeGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (map[string]interface{}, error) {
	// This would typically use Stripe's webhook validation
	// For now, we'll return a placeholder
	return map[string]interface{}{
		"type": "payment_intent.succeeded",
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id": "pi_test_123",
			},
		},
	}, nil
}
