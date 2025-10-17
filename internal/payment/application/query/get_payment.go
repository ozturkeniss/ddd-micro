package query

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// GetPaymentQuery represents the query to get a payment
type GetPaymentQuery struct {
	UserID    uint
	PaymentID string
}

// GetPaymentQueryHandler handles the get payment query
type GetPaymentQueryHandler struct {
	paymentRepo domain.PaymentRepository
}

// NewGetPaymentQueryHandler creates a new get payment query handler
func NewGetPaymentQueryHandler(paymentRepo domain.PaymentRepository) *GetPaymentQueryHandler {
	return &GetPaymentQueryHandler{
		paymentRepo: paymentRepo,
	}
}

// Handle handles the get payment query
func (h *GetPaymentQueryHandler) Handle(ctx context.Context, query GetPaymentQuery) (*dto.PaymentResponse, error) {
	payment, err := h.paymentRepo.GetByID(ctx, query.PaymentID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the payment
	if payment.UserID != query.UserID {
		return nil, domain.ErrPaymentNotFound
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
