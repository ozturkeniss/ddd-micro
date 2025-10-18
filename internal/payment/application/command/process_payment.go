package command

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// ProcessPaymentCommand represents the command to process a payment
type ProcessPaymentCommand struct {
	PaymentID        string
	PaymentMethodID  string
	ConfirmationData map[string]interface{}
}

// ProcessPaymentCommandHandler handles the process payment command
type ProcessPaymentCommandHandler struct {
	paymentRepo    domain.PaymentRepository
	paymentGateway domain.PaymentGateway
}

// NewProcessPaymentCommandHandler creates a new process payment command handler
func NewProcessPaymentCommandHandler(
	paymentRepo domain.PaymentRepository,
	paymentGateway domain.PaymentGateway,
) *ProcessPaymentCommandHandler {
	return &ProcessPaymentCommandHandler{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
	}
}

// Handle handles the process payment command
func (h *ProcessPaymentCommandHandler) Handle(ctx context.Context, cmd ProcessPaymentCommand) (*dto.PaymentResponse, error) {
	// Get payment from repository
	payment, err := h.paymentRepo.GetByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, err
	}

	// Process payment via gateway
	gatewayResponse, err := h.paymentGateway.ProcessPayment(ctx, payment, cmd.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	// Update payment status
	payment.Status = gatewayResponse.Status
	payment.TransactionID = &gatewayResponse.TransactionID
	payment.GatewayResponse = gatewayResponse.GatewayResponse

	// Save updated payment
	if err := h.paymentRepo.Update(ctx, payment); err != nil {
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
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
		CompletedAt:     payment.CompletedAt,
		ExpiresAt:       payment.ExpiresAt,
	}, nil
}
