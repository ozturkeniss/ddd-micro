package client

import (
	"context"
	"fmt"

	basketpb "github.com/ddd-micro/api/proto/basket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BasketClient defines the interface for basket service operations
type BasketClient interface {
	GetBasket(ctx context.Context, userID uint) (*basketpb.Basket, error)
	GetBasketItems(ctx context.Context, userID uint) ([]*basketpb.BasketItem, error)
	ValidateBasket(ctx context.Context, userID uint) (*basketpb.Basket, error)
	ClearBasket(ctx context.Context, userID uint) error
	ReserveItems(ctx context.Context, userID uint, items []*basketpb.BasketItem) error
	ReleaseReservation(ctx context.Context, userID uint) error
}

// basketClient implements BasketClient interface
type basketClient struct {
	conn   *grpc.ClientConn
	client basketpb.BasketServiceClient
}

// NewBasketClient creates a new basket client
func NewBasketClient(basketServiceURL string) (BasketClient, error) {
	conn, err := grpc.Dial(basketServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to basket service: %w", err)
	}

	client := basketpb.NewBasketServiceClient(conn)

	return &basketClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetBasket gets user's basket
func (c *basketClient) GetBasket(ctx context.Context, userID uint) (*basketpb.Basket, error) {
	req := &basketpb.GetBasketRequest{
		UserId: uint32(userID),
	}

	resp, err := c.client.GetBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	return resp.Basket, nil
}

// GetBasketItems gets basket items
func (c *basketClient) GetBasketItems(ctx context.Context, userID uint) ([]*basketpb.BasketItem, error) {
	basket, err := c.GetBasket(ctx, userID)
	if err != nil {
		return nil, err
	}

	return basket.Items, nil
}

// ValidateBasket validates basket contents and availability
func (c *basketClient) ValidateBasket(ctx context.Context, userID uint) (*basketpb.Basket, error) {
	req := &basketpb.GetBasketRequest{
		UserId: uint32(userID),
	}

	resp, err := c.client.GetBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate basket: %w", err)
	}

	// Check if basket is empty
	if len(resp.Basket.Items) == 0 {
		return nil, fmt.Errorf("basket is empty")
	}

	// Check if basket has valid items
	for _, item := range resp.Basket.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductId)
		}
		if item.UnitPrice <= 0 {
			return nil, fmt.Errorf("invalid unit price for product %d", item.ProductId)
		}
	}

	return resp.Basket, nil
}

// ClearBasket clears user's basket
func (c *basketClient) ClearBasket(ctx context.Context, userID uint) error {
	req := &basketpb.ClearBasketRequest{
		UserId: uint32(userID),
	}

	_, err := c.client.ClearBasket(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to clear basket: %w", err)
	}

	return nil
}

// ReserveItems reserves items in the basket
func (c *basketClient) ReserveItems(ctx context.Context, userID uint, items []*basketpb.BasketItem) error {
	req := &basketpb.ReserveItemsRequest{
		UserId: uint32(userID),
		Items:  items,
	}

	_, err := c.client.ReserveItems(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to reserve items: %w", err)
	}

	return nil
}

// ReleaseReservation releases reserved items
func (c *basketClient) ReleaseReservation(ctx context.Context, userID uint) error {
	req := &basketpb.ReleaseReservationRequest{
		UserId: uint32(userID),
	}

	_, err := c.client.ReleaseReservation(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to release reservation: %w", err)
	}

	return nil
}

// Close closes the gRPC connection
func (c *basketClient) Close() error {
	return c.conn.Close()
}
