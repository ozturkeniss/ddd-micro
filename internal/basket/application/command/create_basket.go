package command

import (
	"context"
	"time"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
	"github.com/google/uuid"
)

// CreateBasketCommand represents the command to create a basket
type CreateBasketCommand struct {
	UserID uint
}

// CreateBasketCommandHandler handles the CreateBasketCommand
type CreateBasketCommandHandler struct {
	basketRepo domain.BasketRepository
}

// NewCreateBasketCommandHandler creates a new CreateBasketCommandHandler
func NewCreateBasketCommandHandler(basketRepo domain.BasketRepository) *CreateBasketCommandHandler {
	return &CreateBasketCommandHandler{
		basketRepo: basketRepo,
	}
}

// Handle handles the CreateBasketCommand
func (h *CreateBasketCommandHandler) Handle(ctx context.Context, cmd CreateBasketCommand) (*dto.BasketResponse, error) {
	// Check if basket already exists for this user
	exists, err := h.basketRepo.ExistsByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	
	if exists {
		// Return existing basket
		existingBasket, err := h.basketRepo.GetByUserID(ctx, cmd.UserID)
		if err != nil {
			return nil, err
		}
		return h.mapToResponse(existingBasket), nil
	}
	
	// Create new basket
	basket := &domain.Basket{
		ID:        uuid.New().String(),
		UserID:    cmd.UserID,
		Items:     []domain.BasketItem{},
		Total:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Set expiration time (24 hours)
	basket.SetExpiration(24 * time.Hour)
	
	// Validate basket
	if err := basket.Validate(); err != nil {
		return nil, err
	}
	
	// Save basket
	err = h.basketRepo.Create(ctx, basket)
	if err != nil {
		return nil, err
	}
	
	return h.mapToResponse(basket), nil
}

// mapToResponse maps domain.Basket to dto.BasketResponse
func (h *CreateBasketCommandHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
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
