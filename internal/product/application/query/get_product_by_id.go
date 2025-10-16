package query

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// GetProductByIDQuery represents the query to get a product by ID
type GetProductByIDQuery struct {
	ProductID uint `json:"product_id"`
}

// GetProductByIDHandler handles the get product by ID query
type GetProductByIDHandler struct {
	repo domain.ProductRepository
}

// NewGetProductByIDHandler creates a new get product by ID handler
func NewGetProductByIDHandler(repo domain.ProductRepository) *GetProductByIDHandler {
	return &GetProductByIDHandler{
		repo: repo,
	}
}

// Handle executes the get product by ID query
func (h *GetProductByIDHandler) Handle(ctx context.Context, q GetProductByIDQuery) (*domain.Product, error) {
	return h.repo.GetByID(ctx, q.ProductID)
}

// GetProductBySKUQuery represents the query to get a product by SKU
type GetProductBySKUQuery struct {
	SKU string `json:"sku"`
}

// GetProductBySKUHandler handles the get product by SKU query
type GetProductBySKUHandler struct {
	repo domain.ProductRepository
}

// NewGetProductBySKUHandler creates a new get product by SKU handler
func NewGetProductBySKUHandler(repo domain.ProductRepository) *GetProductBySKUHandler {
	return &GetProductBySKUHandler{
		repo: repo,
	}
}

// Handle executes the get product by SKU query
func (h *GetProductBySKUHandler) Handle(ctx context.Context, q GetProductBySKUQuery) (*domain.Product, error) {
	return h.repo.GetBySKU(ctx, q.SKU)
}
