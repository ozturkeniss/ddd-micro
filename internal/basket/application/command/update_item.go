package command

import (
	"context"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
)

// UpdateItemCommand represents the command to update an item quantity
type UpdateItemCommand struct {
	UserID    uint
	ProductID uint
	Quantity  int
}

// UpdateItemCommandHandler handles the UpdateItemCommand
type UpdateItemCommandHandler struct {
	basketRepo domain.BasketRepository
}

// NewUpdateItemCommandHandler creates a new UpdateItemCommandHandler
func NewUpdateItemCommandHandler(basketRepo domain.BasketRepository) *UpdateItemCommandHandler {
	return &UpdateItemCommandHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the UpdateItemCommand
func (h *UpdateItemCommandHandler) Handle(ctx context.Context, cmd UpdateItemCommand) (*dto.BasketResponse, error) {
	// Get basket for user
	basket, err := h.basketRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Check if basket is expired
	if basket.IsExpired() {
		return nil, domain.ErrBasketExpired
	}

	// Create basket item with updated quantity
	item := &domain.BasketItem{
		BasketID:  basket.ID,
		ProductID: cmd.ProductID,
		Quantity:  cmd.Quantity,
	}

	// Update item in basket
	err = h.basketRepo.UpdateItem(ctx, basket.ID, item)
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
func (h *UpdateItemCommandHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
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
