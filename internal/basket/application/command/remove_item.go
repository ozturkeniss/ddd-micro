package command

import (
	"context"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
)

// RemoveItemCommand represents the command to remove an item from the basket
type RemoveItemCommand struct {
	UserID    uint
	ProductID uint
}

// RemoveItemCommandHandler handles the RemoveItemCommand
type RemoveItemCommandHandler struct {
	basketRepo domain.BasketRepository
}

// NewRemoveItemCommandHandler creates a new RemoveItemCommandHandler
func NewRemoveItemCommandHandler(basketRepo domain.BasketRepository) *RemoveItemCommandHandler {
	return &RemoveItemCommandHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the RemoveItemCommand
func (h *RemoveItemCommandHandler) Handle(ctx context.Context, cmd RemoveItemCommand) (*dto.BasketResponse, error) {
	// Get basket for user
	basket, err := h.basketRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	
	// Check if basket is expired
	if basket.IsExpired() {
		return nil, domain.ErrBasketExpired
	}
	
	// Remove item from basket
	err = h.basketRepo.RemoveItem(ctx, basket.ID, cmd.ProductID)
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
func (h *RemoveItemCommandHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
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
