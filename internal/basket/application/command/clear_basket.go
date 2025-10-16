package command

import (
	"context"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
)

// ClearBasketCommand represents the command to clear all items from the basket
type ClearBasketCommand struct {
	UserID uint
}

// ClearBasketCommandHandler handles the ClearBasketCommand
type ClearBasketCommandHandler struct {
	basketRepo domain.BasketRepository
}

// NewClearBasketCommandHandler creates a new ClearBasketCommandHandler
func NewClearBasketCommandHandler(basketRepo domain.BasketRepository) *ClearBasketCommandHandler {
	return &ClearBasketCommandHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the ClearBasketCommand
func (h *ClearBasketCommandHandler) Handle(ctx context.Context, cmd ClearBasketCommand) (*dto.BasketResponse, error) {
	// Get basket for user
	basket, err := h.basketRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	
	// Check if basket is expired
	if basket.IsExpired() {
		return nil, domain.ErrBasketExpired
	}
	
	// Clear all items from basket
	err = h.basketRepo.ClearItems(ctx, basket.ID)
	if err != nil {
		return nil, err
	}
	
	// Get updated basket
	updatedBasket, err := h.basketRepo.GetByID(ctx, basket.ID)
	if err != nil {
		return nil, err
	}
	
	return h.mapToResponse(updatedBasket), nil
}

// mapToResponse maps domain.Basket to application.BasketResponse
func (h *ClearBasketCommandHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
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
