package command

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// UpdatePaymentMethodCommand represents the command to update a payment method
type UpdatePaymentMethodCommand struct {
	PaymentMethodID string
	IsDefault       bool
	IsActive        bool
}

// UpdatePaymentMethodCommandHandler handles the update payment method command
type UpdatePaymentMethodCommandHandler struct {
	paymentMethodRepo domain.PaymentMethodRepository
}

// NewUpdatePaymentMethodCommandHandler creates a new update payment method command handler
func NewUpdatePaymentMethodCommandHandler(
	paymentMethodRepo domain.PaymentMethodRepository,
) *UpdatePaymentMethodCommandHandler {
	return &UpdatePaymentMethodCommandHandler{
		paymentMethodRepo: paymentMethodRepo,
	}
}

// Handle handles the update payment method command
func (h *UpdatePaymentMethodCommandHandler) Handle(ctx context.Context, cmd UpdatePaymentMethodCommand) (*dto.PaymentMethodResponse, error) {
	// Get payment method from repository
	paymentMethod, err := h.paymentMethodRepo.GetByID(ctx, cmd.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	// Update payment method
	paymentMethod.IsDefault = cmd.IsDefault
	paymentMethod.IsActive = cmd.IsActive

	// Save updated payment method
	if err := h.paymentMethodRepo.Update(ctx, paymentMethod); err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.PaymentMethodResponse{
		ID:             paymentMethod.ID,
		UserID:         paymentMethod.UserID,
		Type:           paymentMethod.Type,
		Provider:       paymentMethod.Provider,
		LastFourDigits: paymentMethod.LastFourDigits,
		ExpiryMonth:    paymentMethod.ExpiryMonth,
		ExpiryYear:     paymentMethod.ExpiryYear,
		IsDefault:      paymentMethod.IsDefault,
		IsActive:       paymentMethod.IsActive,
		CreatedAt:      paymentMethod.CreatedAt,
		UpdatedAt:      paymentMethod.UpdatedAt,
	}, nil
}
