package client

import (
	"context"
	"fmt"

	"github.com/ddd-micro/api/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClient interface for product service operations
type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*product.Product, error)
	ValidateProduct(ctx context.Context, productID uint) error
	CheckStock(ctx context.Context, productID uint, quantity int) error
	Close() error
}

// productClient implements ProductClient interface
type productClient struct {
	conn   *grpc.ClientConn
	client product.ProductServiceClient
}

// NewProductClient creates a new product service gRPC client
func NewProductClient(productServiceURL string) (ProductClient, error) {
	// Create gRPC connection
	conn, err := grpc.Dial(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	// Create client
	client := product.NewProductServiceClient(conn)

	return &productClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetProduct retrieves product information by ID
func (c *productClient) GetProduct(ctx context.Context, productID uint) (*product.Product, error) {
	req := &product.GetProductRequest{
		Id: uint32(productID),
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp.Product, nil
}

// ValidateProduct validates if a product exists and is active
func (c *productClient) ValidateProduct(ctx context.Context, productID uint) error {
	req := &product.GetProductRequest{
		Id: uint32(productID),
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to validate product: %w", err)
	}

	if resp.Product == nil {
		return fmt.Errorf("product not found")
	}

	// Check if product is active
	if resp.Product.IsActive == nil || !*resp.Product.IsActive {
		return fmt.Errorf("product is not active")
	}

	return nil
}

// CheckStock checks if there's enough stock for the requested quantity
func (c *productClient) CheckStock(ctx context.Context, productID uint, quantity int) error {
	req := &product.GetProductRequest{
		Id: uint32(productID),
	}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to check stock: %w", err)
	}

	if resp.Product == nil {
		return fmt.Errorf("product not found")
	}

	// Check if product has stock information
	if resp.Product.Stock == nil {
		return fmt.Errorf("product stock information not available")
	}

	// Check if there's enough stock
	if *resp.Product.Stock < int32(quantity) {
		return fmt.Errorf("insufficient stock: requested %d, available %d", quantity, *resp.Product.Stock)
	}

	return nil
}

// Close closes the gRPC connection
func (c *productClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// HealthCheck performs a health check on the product service
func (c *productClient) HealthCheck(ctx context.Context) error {
	req := &product.HealthCheckRequest{}
	_, err := c.client.HealthCheck(ctx, req)
	if err != nil {
		return fmt.Errorf("product service health check failed: %w", err)
	}
	return nil
}
