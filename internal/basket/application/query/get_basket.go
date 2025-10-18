package query

import (
	"context"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
)

// GetBasketQuery represents the query to get a basket
type GetBasketQuery struct {
	UserID uint
}

// GetBasketQueryHandler handles the GetBasketQuery
type GetBasketQueryHandler struct {
	basketRepo domain.BasketRepository
}

// NewGetBasketQueryHandler creates a new GetBasketQueryHandler
func NewGetBasketQueryHandler(basketRepo domain.BasketRepository) *GetBasketQueryHandler {
	return &GetBasketQueryHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the GetBasketQuery
func (h *GetBasketQueryHandler) Handle(ctx context.Context, query GetBasketQuery) (*dto.BasketResponse, error) {
	// Get basket for user
	basket, err := h.basketRepo.GetByUserID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return h.mapToResponse(basket), nil
}

// mapToResponse maps domain.Basket to application.BasketResponse
func (h *GetBasketQueryHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
	items := make([]dto.BasketItemResponse, len(basket.Items))
	for i, item := range basket.Items {
		items[i] = dto.BasketItemResponse{
			ID:         item.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}
	}

	return &dto.BasketResponse{
		ID:        basket.ID,
		UserID:    basket.UserID,
		Items:     items,
		Total:     basket.Total,
		ItemCount: basket.GetItemCount(),
		CreatedAt: basket.CreatedAt,
		UpdatedAt: basket.UpdatedAt,
		ExpiresAt: basket.ExpiresAt,
		IsExpired: basket.IsExpired(),
	}
}
