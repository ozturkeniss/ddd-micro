package command

import (
	"context"
	"time"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
	"github.com/google/uuid"
)

// CreatePaymentCommand represents the command to create a payment
type CreatePaymentCommand struct {
	UserID          uint
	OrderID         string
	Amount          float64
	Currency        string
	PaymentMethod   string
	PaymentMethodID string
	ReturnURL       string
	CancelURL       string
}

// CreatePaymentCommandHandler handles the create payment command
type CreatePaymentCommandHandler struct {
	paymentRepo    domain.PaymentRepository
	paymentGateway domain.PaymentGateway
}

// NewCreatePaymentCommandHandler creates a new create payment command handler
func NewCreatePaymentCommandHandler(
	paymentRepo domain.PaymentRepository,
	paymentGateway domain.PaymentGateway,
) *CreatePaymentCommandHandler {
	return &CreatePaymentCommandHandler{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
	}
}

// Handle handles the create payment command
func (h *CreatePaymentCommandHandler) Handle(ctx context.Context, cmd CreatePaymentCommand) (*dto.PaymentResponse, error) {
	// Create payment domain object
	payment := &domain.Payment{
		ID:              uuid.New().String(),
		UserID:          cmd.UserID,
		OrderID:         cmd.OrderID,
		Amount:          cmd.Amount,
		Currency:        cmd.Currency,
		Status:          domain.PaymentStatusPending,
		PaymentMethod:   domain.PaymentMethod(cmd.PaymentMethod),
		PaymentProvider: "stripe", // Default provider
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Set expiration time (24 hours)
	payment.SetExpiration(24 * time.Hour)

	// Validate payment
	if err := payment.Validate(); err != nil {
		return nil, err
	}

	// Create payment via gateway
	gatewayResponse, err := h.paymentGateway.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Update payment with gateway response
	payment.TransactionID = &gatewayResponse.TransactionID
	payment.GatewayResponse = gatewayResponse.GatewayResponse

	// Save payment to repository
	if err := h.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.PaymentResponse{
		ID:              payment.ID,
		UserID:          payment.UserID,
		OrderID:         payment.OrderID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Status:          string(payment.Status),
		PaymentMethod:   string(payment.PaymentMethod),
		PaymentProvider: payment.PaymentProvider,
		TransactionID:   payment.TransactionID,
		GatewayResponse: payment.GatewayResponse,
		PaymentURL:      gatewayResponse.PaymentURL,
		ClientSecret:    gatewayResponse.ClientSecret,
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
		CompletedAt:     payment.CompletedAt,
		ExpiresAt:       payment.ExpiresAt,
	}, nil
}
