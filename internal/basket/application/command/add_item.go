package command

import (
	"context"
	"time"

	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/domain"
)

// AddItemCommand represents the command to add an item to the basket
type AddItemCommand struct {
	UserID    uint
	ProductID uint
	Quantity  int
	UnitPrice float64
}

// AddItemCommandHandler handles the AddItemCommand
type AddItemCommandHandler struct {
	basketRepo domain.BasketRepository
}

// NewAddItemCommandHandler creates a new AddItemCommandHandler
func NewAddItemCommandHandler(basketRepo domain.BasketRepository) *AddItemCommandHandler {
	return &AddItemCommandHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the AddItemCommand
func (h *AddItemCommandHandler) Handle(ctx context.Context, cmd AddItemCommand) (*application.BasketResponse, error) {
	// Get or create basket for user
	basket, err := h.basketRepo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		if err == domain.ErrBasketNotFound {
			// Create new basket
			basket = &domain.Basket{
				UserID: cmd.UserID,
				Items:  []domain.BasketItem{},
				Total:  0,
			}
			basket.SetExpiration(24 * time.Hour)
			
			if err := h.basketRepo.Create(ctx, basket); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	// Check if basket is expired
	if basket.IsExpired() {
		return nil, domain.ErrBasketExpired
	}
	
	// Create basket item
	item := &domain.BasketItem{
		BasketID:   basket.ID,
		ProductID:  cmd.ProductID,
		Quantity:   cmd.Quantity,
		UnitPrice:  cmd.UnitPrice,
		TotalPrice: float64(cmd.Quantity) * cmd.UnitPrice,
	}
	
	// Validate item
	if err := item.Validate(); err != nil {
		return nil, err
	}
	
	// Add item to basket
	err = h.basketRepo.AddItem(ctx, basket.ID, item)
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
func (h *AddItemCommandHandler) mapToResponse(basket *domain.Basket) *application.BasketResponse {
	items := make([]application.BasketItemResponse, len(basket.Items))
	for i, item := range basket.Items {
		items[i] = application.BasketItemResponse{
			ID:         item.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}
	}
	
	return &application.BasketResponse{
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
