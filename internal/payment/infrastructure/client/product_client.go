package client

import (
	"context"
	"fmt"

	productpb "github.com/ddd-micro/api/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClient defines the interface for product service operations
type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*productpb.Product, error)
	GetProducts(ctx context.Context, productIDs []uint) ([]*productpb.Product, error)
	ValidateProducts(ctx context.Context, productIDs []uint) ([]*productpb.Product, error)
	UpdateStock(ctx context.Context, productID uint, quantity int) error
}

// productClient implements ProductClient interface
type productClient struct {
	conn   *grpc.ClientConn
	client productpb.ProductServiceClient
}

// NewProductClient creates a new product client
func NewProductClient(productServiceURL string) (ProductClient, error) {
	conn, err := grpc.Dial(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := productpb.NewProductServiceClient(conn)

	return &productClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetProduct gets a single product by ID
func (c *productClient) GetProduct(ctx context.Context, productID uint) (*productpb.Product, error) {
	req := &productpb.GetProductRequest{
		Id: uint32(productID),
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp.Product, nil
}

// GetProducts gets multiple products by IDs
func (c *productClient) GetProducts(ctx context.Context, productIDs []uint) ([]*productpb.Product, error) {
	// Convert uint to uint32
	ids := make([]uint32, len(productIDs))
	for i, id := range productIDs {
		ids[i] = uint32(id)
	}

	// Since there's no GetProducts method, we'll use individual GetProduct calls
	var products []*productpb.Product
	for _, id := range productIDs {
		req := &productpb.GetProductRequest{
			ProductId: uint32(id),
		}

		resp, err := c.client.GetProduct(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", id, err)
		}

		products = append(products, resp)
	}

	return products, nil
}

// ValidateProducts validates that products exist and are available
func (c *productClient) ValidateProducts(ctx context.Context, productIDs []uint) ([]*productpb.Product, error) {
	// Convert uint to uint32
	ids := make([]uint32, len(productIDs))
	for i, id := range productIDs {
		ids[i] = uint32(id)
	}

	// Since there's no ValidateProducts method, we'll use individual GetProduct calls
	var products []*productpb.Product
	for _, id := range productIDs {
		req := &productpb.GetProductRequest{
			ProductId: uint32(id),
		}

		resp, err := c.client.GetProduct(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to validate product %d: %w", id, err)
		}

		// Check if product is available
		if !resp.IsActive {
			return nil, fmt.Errorf("product %d is not active", id)
		}

		products = append(products, resp)
	}

	return products, nil
}

// UpdateStock updates product stock
func (c *productClient) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	req := &productpb.UpdateStockRequest{
		ProductId: uint32(productID),
		Stock:     int32(quantity),
	}

	_, err := c.client.UpdateStock(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

// Close closes the gRPC connection
func (c *productClient) Close() error {
	return c.conn.Close()
}
