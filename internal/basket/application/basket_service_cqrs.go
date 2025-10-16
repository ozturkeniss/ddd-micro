package application

import (
	"context"

	"github.com/ddd-micro/internal/basket/application/command"
	"github.com/ddd-micro/internal/basket/application/dto"
	"github.com/ddd-micro/internal/basket/application/query"
	"github.com/ddd-micro/internal/basket/domain"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
)

// BasketServiceCQRS represents the main basket service using CQRS pattern
type BasketServiceCQRS struct {
	// Command handlers
	createBasketHandler  *command.CreateBasketCommandHandler
	addItemHandler       *command.AddItemCommandHandler
	updateItemHandler    *command.UpdateItemCommandHandler
	removeItemHandler    *command.RemoveItemCommandHandler
	clearBasketHandler   *command.ClearBasketCommandHandler
	
	// Query handlers
	getBasketHandler     *query.GetBasketQueryHandler
	
	// Repository
	basketRepo           domain.BasketRepository
}

// NewBasketServiceCQRS creates a new BasketServiceCQRS
func NewBasketServiceCQRS(basketRepo domain.BasketRepository, productClient client.ProductClient) *BasketServiceCQRS {
	return &BasketServiceCQRS{
		createBasketHandler:  command.NewCreateBasketCommandHandler(basketRepo),
		addItemHandler:       command.NewAddItemCommandHandler(basketRepo, productClient),
		updateItemHandler:    command.NewUpdateItemCommandHandler(basketRepo),
		removeItemHandler:    command.NewRemoveItemCommandHandler(basketRepo),
		clearBasketHandler:   command.NewClearBasketCommandHandler(basketRepo),
		getBasketHandler:     query.NewGetBasketQueryHandler(basketRepo),
		basketRepo:           basketRepo,
	}
}

// CreateBasket creates a new basket for a user
func (s *BasketServiceCQRS) CreateBasket(ctx context.Context, req dto.CreateBasketRequest) (*dto.BasketResponse, error) {
	cmd := command.CreateBasketCommand{
		UserID: req.UserID,
	}
	
	return s.createBasketHandler.Handle(ctx, cmd)
}

// AddItem adds an item to the basket
func (s *BasketServiceCQRS) AddItem(ctx context.Context, userID uint, req dto.AddItemRequest) (*dto.BasketResponse, error) {
	cmd := command.AddItemCommand{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	
	return s.addItemHandler.Handle(ctx, cmd)
}

// UpdateItem updates the quantity of an item in the basket
func (s *BasketServiceCQRS) UpdateItem(ctx context.Context, userID uint, productID uint, req dto.UpdateItemRequest) (*dto.BasketResponse, error) {
	cmd := command.UpdateItemCommand{
		UserID:    userID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}
	
	return s.updateItemHandler.Handle(ctx, cmd)
}

// RemoveItem removes an item from the basket
func (s *BasketServiceCQRS) RemoveItem(ctx context.Context, userID uint, productID uint) (*dto.BasketResponse, error) {
	cmd := command.RemoveItemCommand{
		UserID:    userID,
		ProductID: productID,
	}
	
	return s.removeItemHandler.Handle(ctx, cmd)
}

// ClearBasket removes all items from the basket
func (s *BasketServiceCQRS) ClearBasket(ctx context.Context, userID uint) (*dto.BasketResponse, error) {
	cmd := command.ClearBasketCommand{
		UserID: userID,
	}
	
	return s.clearBasketHandler.Handle(ctx, cmd)
}

// GetBasket retrieves the basket for a user
func (s *BasketServiceCQRS) GetBasket(ctx context.Context, userID uint) (*dto.BasketResponse, error) {
	query := query.GetBasketQuery{
		UserID: userID,
	}
	
	return s.getBasketHandler.Handle(ctx, query)
}

// DeleteBasket deletes a basket for a user
func (s *BasketServiceCQRS) DeleteBasket(ctx context.Context, userID uint) error {
	return s.basketRepo.DeleteByUserID(ctx, userID)
}

// CleanupExpiredBaskets removes expired baskets
func (s *BasketServiceCQRS) CleanupExpiredBaskets(ctx context.Context) error {
	return s.basketRepo.CleanupExpired(ctx)
}
