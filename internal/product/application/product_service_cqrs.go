package application

import (
	"context"

	"github.com/ddd-micro/internal/product/application/command"
	"github.com/ddd-micro/internal/product/application/query"
	"github.com/ddd-micro/internal/product/domain"
)

// ProductServiceCQRS handles product business logic using CQRS pattern
type ProductServiceCQRS struct {
	// Command handlers
	createProductHandler      *command.CreateProductHandler
	updateProductHandler      *command.UpdateProductHandler
	deleteProductHandler      *command.DeleteProductHandler
	updateStockHandler        *command.UpdateStockHandler
	reduceStockHandler        *command.ReduceStockHandler
	increaseStockHandler      *command.IncreaseStockHandler
	activateProductHandler    *command.ActivateProductHandler
	deactivateProductHandler  *command.DeactivateProductHandler
	markAsFeaturedHandler     *command.MarkAsFeaturedHandler
	unmarkAsFeaturedHandler   *command.UnmarkAsFeaturedHandler
	incrementViewCountHandler *command.IncrementViewCountHandler

	// Query handlers
	getProductByIDHandler         *query.GetProductByIDHandler
	getProductBySKUHandler        *query.GetProductBySKUHandler
	listProductsHandler           *query.ListProductsHandler
	listProductsByCategoryHandler *query.ListProductsByCategoryHandler
	searchProductsHandler         *query.SearchProductsHandler
}

// NewProductServiceCQRS creates a new CQRS-based product service
func NewProductServiceCQRS(repo domain.ProductRepository) *ProductServiceCQRS {
	return &ProductServiceCQRS{
		// Initialize command handlers
		createProductHandler:      command.NewCreateProductHandler(repo),
		updateProductHandler:      command.NewUpdateProductHandler(repo),
		deleteProductHandler:      command.NewDeleteProductHandler(repo),
		updateStockHandler:        command.NewUpdateStockHandler(repo),
		reduceStockHandler:        command.NewReduceStockHandler(repo),
		increaseStockHandler:      command.NewIncreaseStockHandler(repo),
		activateProductHandler:    command.NewActivateProductHandler(repo),
		deactivateProductHandler:  command.NewDeactivateProductHandler(repo),
		markAsFeaturedHandler:     command.NewMarkAsFeaturedHandler(repo),
		unmarkAsFeaturedHandler:   command.NewUnmarkAsFeaturedHandler(repo),
		incrementViewCountHandler: command.NewIncrementViewCountHandler(repo),

		// Initialize query handlers
		getProductByIDHandler:         query.NewGetProductByIDHandler(repo),
		getProductBySKUHandler:        query.NewGetProductBySKUHandler(repo),
		listProductsHandler:           query.NewListProductsHandler(repo),
		listProductsByCategoryHandler: query.NewListProductsByCategoryHandler(repo),
		searchProductsHandler:         query.NewSearchProductsHandler(repo),
	}
}

// ========== COMMAND METHODS ==========

// CreateProduct creates a new product
func (s *ProductServiceCQRS) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	cmd := command.CreateProductCommand{
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
	}

	product, err := s.createProductHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UpdateProduct updates an existing product
func (s *ProductServiceCQRS) UpdateProduct(ctx context.Context, id uint, req UpdateProductRequest) (*ProductResponse, error) {
	cmd := command.UpdateProductCommand{
		ProductID:        id,
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
		Barcode:          req.Barcode,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		Color:            req.Color,
		Size:             req.Size,
		Material:         req.Material,
		Tags:             req.Tags,
		Images:           req.Images,
		IsActive:         req.IsActive,
		IsDigital:        req.IsDigital,
		IsFeatured:       req.IsFeatured,
		IsOnSale:         req.IsOnSale,
		SortOrder:        req.SortOrder,
	}

	product, err := s.updateProductHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct soft deletes a product
func (s *ProductServiceCQRS) DeleteProduct(ctx context.Context, id uint) error {
	cmd := command.DeleteProductCommand{
		ProductID: id,
	}

	return s.deleteProductHandler.Handle(ctx, cmd)
}

// UpdateStock updates the stock of a product
func (s *ProductServiceCQRS) UpdateStock(ctx context.Context, id uint, stock int) error {
	cmd := command.UpdateStockCommand{
		ProductID: id,
		Stock:     stock,
	}

	return s.updateStockHandler.Handle(ctx, cmd)
}

// ReduceStock reduces the stock of a product
func (s *ProductServiceCQRS) ReduceStock(ctx context.Context, id uint, amount int) error {
	cmd := command.ReduceStockCommand{
		ProductID: id,
		Amount:    amount,
	}

	return s.reduceStockHandler.Handle(ctx, cmd)
}

// IncreaseStock increases the stock of a product
func (s *ProductServiceCQRS) IncreaseStock(ctx context.Context, id uint, amount int) error {
	cmd := command.IncreaseStockCommand{
		ProductID: id,
		Amount:    amount,
	}

	return s.increaseStockHandler.Handle(ctx, cmd)
}

// ActivateProduct activates a product
func (s *ProductServiceCQRS) ActivateProduct(ctx context.Context, id uint) (*ProductResponse, error) {
	cmd := command.ActivateProductCommand{
		ProductID: id,
	}

	product, err := s.activateProductHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeactivateProduct deactivates a product
func (s *ProductServiceCQRS) DeactivateProduct(ctx context.Context, id uint) (*ProductResponse, error) {
	cmd := command.DeactivateProductCommand{
		ProductID: id,
	}

	product, err := s.deactivateProductHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// MarkAsFeatured marks a product as featured
func (s *ProductServiceCQRS) MarkAsFeatured(ctx context.Context, id uint) (*ProductResponse, error) {
	cmd := command.MarkAsFeaturedCommand{
		ProductID: id,
	}

	product, err := s.markAsFeaturedHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UnmarkAsFeatured removes featured status from a product
func (s *ProductServiceCQRS) UnmarkAsFeatured(ctx context.Context, id uint) (*ProductResponse, error) {
	cmd := command.UnmarkAsFeaturedCommand{
		ProductID: id,
	}

	product, err := s.unmarkAsFeaturedHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// IncrementViewCount increments the view count of a product
func (s *ProductServiceCQRS) IncrementViewCount(ctx context.Context, id uint) error {
	cmd := command.IncrementViewCountCommand{
		ProductID: id,
	}

	return s.incrementViewCountHandler.Handle(ctx, cmd)
}

// ========== QUERY METHODS ==========

// GetProductByID retrieves a product by ID
func (s *ProductServiceCQRS) GetProductByID(ctx context.Context, id uint) (*ProductResponse, error) {
	q := query.GetProductByIDQuery{
		ProductID: id,
	}

	product, err := s.getProductByIDHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductServiceCQRS) GetProductBySKU(ctx context.Context, sku string) (*ProductResponse, error) {
	q := query.GetProductBySKUQuery{
		SKU: sku,
	}

	product, err := s.getProductBySKUHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// ListProducts retrieves all products with pagination
func (s *ProductServiceCQRS) ListProducts(ctx context.Context, offset, limit int) (*ListProductsResponse, error) {
	q := query.ListProductsQuery{
		Offset: offset,
		Limit:  limit,
	}

	result, err := s.listProductsHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(result.Products))
	for i, product := range result.Products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    result.Total,
		Offset:   result.Offset,
		Limit:    result.Limit,
	}, nil
}

// ListProductsByCategory retrieves products by category with pagination
func (s *ProductServiceCQRS) ListProductsByCategory(ctx context.Context, category string, offset, limit int) (*ListProductsResponse, error) {
	q := query.ListProductsByCategoryQuery{
		Category: category,
		Offset:   offset,
		Limit:    limit,
	}

	result, err := s.listProductsByCategoryHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(result.Products))
	for i, product := range result.Products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    result.Total,
		Offset:   result.Offset,
		Limit:    result.Limit,
	}, nil
}

// SearchProducts searches products by name with pagination
func (s *ProductServiceCQRS) SearchProducts(ctx context.Context, name string, offset, limit int) (*ListProductsResponse, error) {
	q := query.SearchProductsQuery{
		Name:   name,
		Offset: offset,
		Limit:  limit,
	}

	result, err := s.searchProductsHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	productResponses := make([]ProductResponse, len(result.Products))
	for i, product := range result.Products {
		productResponses[i] = *s.toProductResponse(product)
	}

	return &ListProductsResponse{
		Products: productResponses,
		Total:    result.Total,
		Offset:   result.Offset,
		Limit:    result.Limit,
	}, nil
}

// ========== HELPER METHODS ==========

// toProductResponse converts domain.Product to ProductResponse
func (s *ProductServiceCQRS) toProductResponse(product *domain.Product) *ProductResponse {
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
