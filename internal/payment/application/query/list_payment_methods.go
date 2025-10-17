package query

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// ListPaymentMethodsQuery represents the query to list payment methods
type ListPaymentMethodsQuery struct {
	UserID uint
}

// ListPaymentMethodsQueryHandler handles the list payment methods query
type ListPaymentMethodsQueryHandler struct {
	paymentMethodRepo domain.PaymentMethodRepository
}

// NewListPaymentMethodsQueryHandler creates a new list payment methods query handler
func NewListPaymentMethodsQueryHandler(paymentMethodRepo domain.PaymentMethodRepository) *ListPaymentMethodsQueryHandler {
	return &ListPaymentMethodsQueryHandler{
		paymentMethodRepo: paymentMethodRepo,
	}
}

// Handle handles the list payment methods query
func (h *ListPaymentMethodsQueryHandler) Handle(ctx context.Context, query ListPaymentMethodsQuery) (*dto.PaymentMethodListResponse, error) {
	// Get payment methods from repository
	paymentMethods, err := h.paymentMethodRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	paymentMethodDTOs := make([]dto.PaymentMethodResponse, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		paymentMethodDTOs[i] = dto.PaymentMethodResponse{
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
		}
	}

	return &dto.PaymentMethodListResponse{
		PaymentMethods: paymentMethodDTOs,
		Total:          len(paymentMethodDTOs),
	}, nil
}
