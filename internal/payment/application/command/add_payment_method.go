package command

import (
	"context"
	"time"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
	"github.com/google/uuid"
)

// AddPaymentMethodCommand represents the command to add a payment method
type AddPaymentMethodCommand struct {
	UserID    uint
	Type      string
	Provider  string
	Token     string
	IsDefault bool
}

// AddPaymentMethodCommandHandler handles the add payment method command
type AddPaymentMethodCommandHandler struct {
	paymentMethodRepo domain.PaymentMethodRepository
	paymentGateway    domain.PaymentGateway
}

// NewAddPaymentMethodCommandHandler creates a new add payment method command handler
func NewAddPaymentMethodCommandHandler(
	paymentMethodRepo domain.PaymentMethodRepository,
	paymentGateway domain.PaymentGateway,
) *AddPaymentMethodCommandHandler {
	return &AddPaymentMethodCommandHandler{
		paymentMethodRepo: paymentMethodRepo,
		paymentGateway:    paymentGateway,
	}
}

// Handle handles the add payment method command
func (h *AddPaymentMethodCommandHandler) Handle(ctx context.Context, cmd AddPaymentMethodCommand) (*dto.PaymentMethodResponse, error) {
	// Create payment method domain object
	paymentMethod := &domain.PaymentMethodInfo{
		ID:        uuid.New().String(),
		UserID:    cmd.UserID,
		Type:      cmd.Type,
		Provider:  cmd.Provider,
		IsDefault: cmd.IsDefault,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate payment method
	if err := paymentMethod.Validate(); err != nil {
		return nil, err
	}

	// Create payment method via gateway
	_, err := h.paymentGateway.CreatePaymentMethod(ctx, cmd.UserID, paymentMethod)
	if err != nil {
		return nil, err
	}

	// Update payment method with gateway response
	// This would typically include last four digits, expiry info, etc.
	// For now, we'll just save the basic info

	// Save payment method to repository
	if err := h.paymentMethodRepo.Create(ctx, paymentMethod); err != nil {
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
