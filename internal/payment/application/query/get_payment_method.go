package query

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// GetPaymentMethodQuery represents the query to get a payment method
type GetPaymentMethodQuery struct {
	UserID          uint
	PaymentMethodID string
}

// GetPaymentMethodQueryHandler handles the get payment method query
type GetPaymentMethodQueryHandler struct {
	paymentMethodRepo domain.PaymentMethodRepository
}

// NewGetPaymentMethodQueryHandler creates a new get payment method query handler
func NewGetPaymentMethodQueryHandler(paymentMethodRepo domain.PaymentMethodRepository) *GetPaymentMethodQueryHandler {
	return &GetPaymentMethodQueryHandler{
		paymentMethodRepo: paymentMethodRepo,
	}
}

// Handle handles the get payment method query
func (h *GetPaymentMethodQueryHandler) Handle(ctx context.Context, query GetPaymentMethodQuery) (*dto.PaymentMethodResponse, error) {
	paymentMethod, err := h.paymentMethodRepo.GetByID(ctx, query.PaymentMethodID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the payment method
	if paymentMethod.UserID != query.UserID {
		return nil, domain.ErrPaymentMethodNotFound
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
