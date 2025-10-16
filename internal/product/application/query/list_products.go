package query

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// ListProductsQuery represents the query to list products
type ListProductsQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// ListProductsResult represents the result of listing products
type ListProductsResult struct {
	Products []*domain.Product `json:"products"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

// ListProductsHandler handles the list products query
type ListProductsHandler struct {
	repo domain.ProductRepository
}

// NewListProductsHandler creates a new list products handler
func NewListProductsHandler(repo domain.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{
		repo: repo,
	}
}

// Handle executes the list products query
func (h *ListProductsHandler) Handle(ctx context.Context, q ListProductsQuery) (*ListProductsResult, error) {
	products, err := h.repo.List(ctx, q.Offset, q.Limit)
	if err != nil {
		return nil, err
	}

	return &ListProductsResult{
		Products: products,
		Total:    len(products),
		Offset:   q.Offset,
		Limit:    q.Limit,
	}, nil
}

// ListProductsByCategoryQuery represents the query to list products by category
type ListProductsByCategoryQuery struct {
	Category string `json:"category"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

// ListProductsByCategoryResult represents the result of listing products by category
type ListProductsByCategoryResult struct {
	Products []*domain.Product `json:"products"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

// ListProductsByCategoryHandler handles the list products by category query
type ListProductsByCategoryHandler struct {
	repo domain.ProductRepository
}

// NewListProductsByCategoryHandler creates a new list products by category handler
func NewListProductsByCategoryHandler(repo domain.ProductRepository) *ListProductsByCategoryHandler {
	return &ListProductsByCategoryHandler{
		repo: repo,
	}
}

// Handle executes the list products by category query
func (h *ListProductsByCategoryHandler) Handle(ctx context.Context, q ListProductsByCategoryQuery) (*ListProductsByCategoryResult, error) {
	products, err := h.repo.ListByCategory(ctx, q.Category, q.Offset, q.Limit)
	if err != nil {
		return nil, err
	}

	return &ListProductsByCategoryResult{
		Products: products,
		Total:    len(products),
		Offset:   q.Offset,
		Limit:    q.Limit,
	}, nil
}

// SearchProductsQuery represents the query to search products
type SearchProductsQuery struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// SearchProductsResult represents the result of searching products
type SearchProductsResult struct {
	Products []*domain.Product `json:"products"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

// SearchProductsHandler handles the search products query
type SearchProductsHandler struct {
	repo domain.ProductRepository
}

// NewSearchProductsHandler creates a new search products handler
func NewSearchProductsHandler(repo domain.ProductRepository) *SearchProductsHandler {
	return &SearchProductsHandler{
		repo: repo,
	}
}

// Handle executes the search products query
func (h *SearchProductsHandler) Handle(ctx context.Context, q SearchProductsQuery) (*SearchProductsResult, error) {
	products, err := h.repo.SearchByName(ctx, q.Name, q.Offset, q.Limit)
	if err != nil {
		return nil, err
	}

	return &SearchProductsResult{
		Products: products,
		Total:    len(products),
		Offset:   q.Offset,
		Limit:    q.Limit,
	}, nil
}
