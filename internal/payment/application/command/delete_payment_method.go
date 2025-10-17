package command

import (
	"context"

	"github.com/ddd-micro/internal/payment/domain"
)

// DeletePaymentMethodCommand represents the command to delete a payment method
type DeletePaymentMethodCommand struct {
	PaymentMethodID string
}

// DeletePaymentMethodCommandHandler handles the delete payment method command
type DeletePaymentMethodCommandHandler struct {
	paymentMethodRepo domain.PaymentMethodRepository
}

// NewDeletePaymentMethodCommandHandler creates a new delete payment method command handler
func NewDeletePaymentMethodCommandHandler(
	paymentMethodRepo domain.PaymentMethodRepository,
) *DeletePaymentMethodCommandHandler {
	return &DeletePaymentMethodCommandHandler{
		paymentMethodRepo: paymentMethodRepo,
	}
}

// Handle handles the delete payment method command
func (h *DeletePaymentMethodCommandHandler) Handle(ctx context.Context, cmd DeletePaymentMethodCommand) error {
	// Delete payment method from repository
	return h.paymentMethodRepo.Delete(ctx, cmd.PaymentMethodID)
}
