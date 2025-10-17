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
func NewBasketServiceCQRS(basketRepo domain.BasketRepository, userClient client.UserClient, productClient client.ProductClient) *BasketServiceCQRS {
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

// AddItem adds an item to the basket (HTTP version)
func (s *BasketServiceCQRS) AddItemHTTP(ctx context.Context, userID uint, req dto.AddItemRequest) (*dto.BasketResponse, error) {
	cmd := command.AddItemCommand{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	
	return s.addItemHandler.Handle(ctx, cmd)
}

// UpdateItem updates the quantity of an item in the basket (HTTP version)
func (s *BasketServiceCQRS) UpdateItemHTTP(ctx context.Context, userID uint, productID uint, req dto.UpdateItemRequest) (*dto.BasketResponse, error) {
	cmd := command.UpdateItemCommand{
		UserID:    userID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}
	
	return s.updateItemHandler.Handle(ctx, cmd)
}

// RemoveItem removes an item from the basket (HTTP version)
func (s *BasketServiceCQRS) RemoveItemHTTP(ctx context.Context, userID uint, productID uint) (*dto.BasketResponse, error) {
	cmd := command.RemoveItemCommand{
		UserID:    userID,
		ProductID: productID,
	}
	
	return s.removeItemHandler.Handle(ctx, cmd)
}

// ClearBasket removes all items from the basket (HTTP version)
func (s *BasketServiceCQRS) ClearBasketHTTP(ctx context.Context, userID uint) (*dto.BasketResponse, error) {
	cmd := command.ClearBasketCommand{
		UserID: userID,
	}
	
	return s.clearBasketHandler.Handle(ctx, cmd)
}

// GetBasket retrieves the basket for a user (HTTP version)
func (s *BasketServiceCQRS) GetBasketHTTP(ctx context.Context, userID uint) (*dto.BasketResponse, error) {
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
	_, err := s.basketRepo.CleanupExpired(ctx)
	return err
}

// AddItem adds an item to the basket (gRPC version)
func (s *BasketServiceCQRS) AddItem(ctx context.Context, req dto.AddItemRequest) (*dto.BasketResponse, error) {
	cmd := command.AddItemCommand{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}
	
	return s.addItemHandler.Handle(ctx, cmd)
}

// UpdateItem updates the quantity of an item in the basket (gRPC version)
func (s *BasketServiceCQRS) UpdateItem(ctx context.Context, productID uint, req dto.UpdateItemRequest) (*dto.BasketResponse, error) {
	cmd := command.UpdateItemCommand{
		UserID:    req.UserID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}
	
	return s.updateItemHandler.Handle(ctx, cmd)
}

// RemoveItem removes an item from the basket (gRPC version)
func (s *BasketServiceCQRS) RemoveItem(ctx context.Context, req dto.RemoveItemRequest) (*dto.BasketResponse, error) {
	cmd := command.RemoveItemCommand{
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}
	
	return s.removeItemHandler.Handle(ctx, cmd)
}

// ClearBasket removes all items from the basket (gRPC version)
func (s *BasketServiceCQRS) ClearBasket(ctx context.Context, req dto.ClearBasketRequest) error {
	cmd := command.ClearBasketCommand{
		UserID: req.UserID,
	}
	
	_, err := s.clearBasketHandler.Handle(ctx, cmd)
	return err
}

// GetBasket retrieves the basket for a user (gRPC version)
func (s *BasketServiceCQRS) GetBasket(ctx context.Context, req dto.GetBasketRequest) (*dto.BasketResponse, error) {
	query := query.GetBasketQuery{
		UserID: req.UserID,
	}
	
	return s.getBasketHandler.Handle(ctx, query)
}

// AdminCleanupExpiredBaskets removes expired baskets and returns count
func (s *BasketServiceCQRS) AdminCleanupExpiredBaskets(ctx context.Context) (int, error) {
	return s.basketRepo.CleanupExpired(ctx)
}
