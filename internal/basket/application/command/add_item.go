package command

import (
	"context"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/domain"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
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
	basketRepo    domain.BasketRepository
	productClient client.ProductClient
}

// NewAddItemCommandHandler creates a new AddItemCommandHandler
func NewAddItemCommandHandler(basketRepo domain.BasketRepository, productClient client.ProductClient) *AddItemCommandHandler {
	return &AddItemCommandHandler{
		basketRepo:    basketRepo,
		productClient: productClient,
	}
}

// Handle handles the AddItemCommand
func (h *AddItemCommandHandler) Handle(ctx context.Context, cmd AddItemCommand) (*dto.BasketResponse, error) {
	// Validate product exists and is active
	if err := h.productClient.ValidateProduct(ctx, cmd.ProductID); err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	// Check stock availability
	if err := h.productClient.CheckStock(ctx, cmd.ProductID, cmd.Quantity); err != nil {
		return nil, fmt.Errorf("stock check failed: %w", err)
	}

	// Get product to get current price
	product, err := h.productClient.GetProduct(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Use current product price if not provided
	unitPrice := cmd.UnitPrice
	if unitPrice == 0 && product.Price != 0 {
		unitPrice = float64(product.Price)
	}

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
		UnitPrice:  unitPrice,
		TotalPrice: float64(cmd.Quantity) * unitPrice,
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
func (h *AddItemCommandHandler) mapToResponse(basket *domain.Basket) *dto.BasketResponse {
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
