package command

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// CancelPaymentCommand represents the command to cancel a payment
type CancelPaymentCommand struct {
	PaymentID string
}

// CancelPaymentCommandHandler handles the cancel payment command
type CancelPaymentCommandHandler struct {
	paymentRepo    domain.PaymentRepository
	paymentGateway domain.PaymentGateway
}

// NewCancelPaymentCommandHandler creates a new cancel payment command handler
func NewCancelPaymentCommandHandler(
	paymentRepo domain.PaymentRepository,
	paymentGateway domain.PaymentGateway,
) *CancelPaymentCommandHandler {
	return &CancelPaymentCommandHandler{
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
	}
}

// Handle handles the cancel payment command
func (h *CancelPaymentCommandHandler) Handle(ctx context.Context, cmd CancelPaymentCommand) (*dto.PaymentResponse, error) {
	// Get payment from repository
	payment, err := h.paymentRepo.GetByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, err
	}

	// Check if payment can be cancelled
	if !payment.CanBeCancelled() {
		return nil, domain.ErrPaymentCannotBeCancelled
	}

	// Cancel payment via gateway
	gatewayResponse, err := h.paymentGateway.CancelPayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Update payment status
	payment.SetCancelled()
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
