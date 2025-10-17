package query

import (
	"context"

	"github.com/ddd-micro/internal/payment/application/dto"
	"github.com/ddd-micro/internal/payment/domain"
)

// ListPaymentsQuery represents the query to list payments
type ListPaymentsQuery struct {
	UserID uint
	Page   int
	Limit  int
	Status string
}

// ListPaymentsQueryHandler handles the list payments query
type ListPaymentsQueryHandler struct {
	paymentRepo domain.PaymentRepository
}

// NewListPaymentsQueryHandler creates a new list payments query handler
func NewListPaymentsQueryHandler(paymentRepo domain.PaymentRepository) *ListPaymentsQueryHandler {
	return &ListPaymentsQueryHandler{
		paymentRepo: paymentRepo,
	}
}

// Handle handles the list payments query
func (h *ListPaymentsQueryHandler) Handle(ctx context.Context, query ListPaymentsQuery) (*dto.PaymentListResponse, error) {
	// Calculate offset
	offset := (query.Page - 1) * query.Limit

	// Get payments from repository
	payments, total, err := h.paymentRepo.GetByUserID(ctx, query.UserID, query.Limit, offset, query.Status)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	paymentDTOs := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentDTOs[i] = dto.PaymentResponse{
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
		}
	}

	// Calculate pagination info
	totalPages := (total + query.Limit - 1) / query.Limit
	hasNext := query.Page < totalPages
	hasPrev := query.Page > 1

	return &dto.PaymentListResponse{
		Payments:   paymentDTOs,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}
