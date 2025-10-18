package application

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/product/domain"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product with this SKU already exists")
	ErrInvalidProductData   = errors.New("invalid product data")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidStockAmount   = errors.New("invalid stock amount")
)

// ProductService handles product business logic
type ProductService struct {
	repo domain.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	// Validate product data
	product := &domain.Product{
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		CostPrice:        req.CostPrice,
		Stock:            req.Stock,
		MinStock:         req.MinStock,
		MaxStock:         req.MaxStock,
		Category:         req.Category,
		SubCategory:      req.SubCategory,
		Brand:            req.Brand,
		SKU:              req.SKU,
		Barcode:          req.Barcode,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Color:            req.Color,
		Size:             req.Size,
		Material:         req.Material,
		Tags:             req.Tags,
		Images:           req.Images,
		IsDigital:        req.IsDigital,
		IsFeatured:       req.IsFeatured,
		IsOnSale:         req.IsOnSale,
		SortOrder:        req.SortOrder,
		IsActive:         true,
	}

	// Validate product
	if err := product.ValidateProduct(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*ProductResponse, error) {
	product, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id uint, req UpdateProductRequest) (*ProductResponse, error) {
	// Get existing product
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.ShortDescription != nil {
		product.ShortDescription = *req.ShortDescription
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.ComparePrice != nil {
		product.ComparePrice = *req.ComparePrice
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.MinStock != nil {
		product.MinStock = *req.MinStock
	}
	if req.MaxStock != nil {
		product.MaxStock = *req.MaxStock
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.SubCategory != nil {
		product.SubCategory = *req.SubCategory
	}
	if req.Brand != nil {
		product.Brand = *req.Brand
	}
	if req.Barcode != nil {
		product.Barcode = *req.Barcode
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.Dimensions != nil {
		product.Dimensions = *req.Dimensions
	}
	if req.Color != nil {
		product.Color = *req.Color
	}
	if req.Size != nil {
		product.Size = *req.Size
	}
	if req.Material != nil {
		product.Material = *req.Material
	}
	if req.Tags != nil {
		product.Tags = *req.Tags
	}
	if req.Images != nil {
		product.Images = *req.Images
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.IsOnSale != nil {
		product.IsOnSale = *req.IsOnSale
	}
	if req.SortOrder != nil {
		product.SortOrder = *req.SortOrder
	}

	// Validate updated product
	if err := product.ValidateProduct(); err != nil {
		return nil, err
	}

	// Save changes
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ListProducts retrieves all products with pagination
func (s *ProductService) ListProducts(ctx context.Context, offset, limit int) (*ListProductsResponse, error) {
	products, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Offset:   offset,
		Limit:    limit,
	}, nil
}

// ListProductsByCategory retrieves products by category with pagination
func (s *ProductService) ListProductsByCategory(ctx context.Context, category string, offset, limit int) (*ListProductsResponse, error) {
	products, err := s.repo.ListByCategory(ctx, category, offset, limit)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Offset:   offset,
		Limit:    limit,
	}, nil
}

// SearchProducts searches products by name with pagination
func (s *ProductService) SearchProducts(ctx context.Context, name string, offset, limit int) (*ListProductsResponse, error) {
	products, err := s.repo.SearchByName(ctx, name, offset, limit)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    len(productResponses),
		Offset:   offset,
		Limit:    limit,
	}, nil
}

// UpdateStock updates the stock of a product
func (s *ProductService) UpdateStock(ctx context.Context, id uint, stock int) error {
	return s.repo.UpdateStock(ctx, id, stock)
}

// ReduceStock reduces the stock of a product
func (s *ProductService) ReduceStock(ctx context.Context, id uint, amount int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.ReduceStock(amount); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

// IncreaseStock increases the stock of a product
func (s *ProductService) IncreaseStock(ctx context.Context, id uint, amount int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.IncreaseStock(amount); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

// ActivateProduct activates a product
func (s *ProductService) ActivateProduct(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Activate()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeactivateProduct deactivates a product
func (s *ProductService) DeactivateProduct(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Deactivate()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// MarkAsFeatured marks a product as featured
func (s *ProductService) MarkAsFeatured(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.MarkAsFeatured()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UnmarkAsFeatured removes featured status from a product
func (s *ProductService) UnmarkAsFeatured(ctx context.Context, id uint) (*ProductResponse, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.UnmarkAsFeatured()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// IncrementViewCount increments the view count of a product
func (s *ProductService) IncrementViewCount(ctx context.Context, id uint) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	product.IncrementViewCount()

	return s.repo.Update(ctx, product)
}

// toProductResponse converts domain.Product to ProductResponse
func (s *ProductService) toProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:               product.ID,
		Name:             product.Name,
		Description:      product.Description,
		ShortDescription: product.ShortDescription,
		Price:            product.Price,
		ComparePrice:     product.ComparePrice,
		CostPrice:        product.CostPrice,
		Stock:            product.Stock,
		MinStock:         product.MinStock,
		MaxStock:         product.MaxStock,
		Category:         product.Category,
		SubCategory:      product.SubCategory,
		Brand:            product.Brand,
		SKU:              product.SKU,
		Barcode:          product.Barcode,
		Weight:           product.Weight,
		Dimensions:       product.Dimensions,
		Color:            product.Color,
		Size:             product.Size,
		Material:         product.Material,
		Tags:             product.Tags,
		Images:           product.Images,
		IsActive:         product.IsActive,
		IsDigital:        product.IsDigital,
		IsFeatured:       product.IsFeatured,
		IsOnSale:         product.IsOnSale,
		SortOrder:        product.SortOrder,
		ViewCount:        product.ViewCount,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}
}
