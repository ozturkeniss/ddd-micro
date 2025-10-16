package application

import (
	"time"
)

// ========== PRODUCT DTOs ==========

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price" binding:"required,min=0"`
	ComparePrice     float64 `json:"compare_price"`
	CostPrice        float64 `json:"cost_price"`
	Stock            int     `json:"stock" binding:"min=0"`
	MinStock         int     `json:"min_stock"`
	MaxStock         int     `json:"max_stock"`
	Category         string  `json:"category"`
	SubCategory      string  `json:"sub_category"`
	Brand            string  `json:"brand"`
	SKU              string  `json:"sku" binding:"required"`
	Barcode          string  `json:"barcode"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
	Color            string  `json:"color"`
	Size             string  `json:"size"`
	Material         string  `json:"material"`
	Tags             string  `json:"tags"`
	Images           string  `json:"images"`
	IsDigital        bool    `json:"is_digital"`
	IsFeatured       bool    `json:"is_featured"`
	IsOnSale         bool    `json:"is_on_sale"`
	SortOrder        int     `json:"sort_order"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name             *string  `json:"name"`
	Description      *string  `json:"description"`
	ShortDescription *string  `json:"short_description"`
	Price            *float64 `json:"price"`
	ComparePrice     *float64 `json:"compare_price"`
	CostPrice        *float64 `json:"cost_price"`
	Stock            *int     `json:"stock"`
	MinStock         *int     `json:"min_stock"`
	MaxStock         *int     `json:"max_stock"`
	Category         *string  `json:"category"`
	SubCategory      *string  `json:"sub_category"`
	Brand            *string  `json:"brand"`
	Barcode          *string  `json:"barcode"`
	Weight           *float64 `json:"weight"`
	Dimensions       *string  `json:"dimensions"`
	Color            *string  `json:"color"`
	Size             *string  `json:"size"`
	Material         *string  `json:"material"`
	Tags             *string  `json:"tags"`
	Images           *string  `json:"images"`
	IsActive         *bool    `json:"is_active"`
	IsDigital        *bool    `json:"is_digital"`
	IsFeatured       *bool    `json:"is_featured"`
	IsOnSale         *bool    `json:"is_on_sale"`
	SortOrder        *int     `json:"sort_order"`
}

// UpdateStockRequest represents the request to update product stock
type UpdateStockRequest struct {
	Stock int `json:"stock" binding:"required,min=0"`
}

// ProductResponse represents the product response
type ProductResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	Price            float64   `json:"price"`
	ComparePrice     float64   `json:"compare_price"`
	CostPrice        float64   `json:"cost_price"`
	Stock            int       `json:"stock"`
	MinStock         int       `json:"min_stock"`
	MaxStock         int       `json:"max_stock"`
	Category         string    `json:"category"`
	SubCategory      string    `json:"sub_category"`
	Brand            string    `json:"brand"`
	SKU              string    `json:"sku"`
	Barcode          string    `json:"barcode"`
	Weight           float64   `json:"weight"`
	Dimensions       string    `json:"dimensions"`
	Color            string    `json:"color"`
	Size             string    `json:"size"`
	Material         string    `json:"material"`
	Tags             string    `json:"tags"`
	Images           string    `json:"images"`
	IsActive         bool      `json:"is_active"`
	IsDigital        bool      `json:"is_digital"`
	IsFeatured       bool      `json:"is_featured"`
	IsOnSale         bool      `json:"is_on_sale"`
	SortOrder        int       `json:"sort_order"`
	ViewCount        int       `json:"view_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ListProductsResponse represents the paginated list of products
type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}

// SearchProductsRequest represents the request to search products
type SearchProductsRequest struct {
	Query    string `json:"query"`
	Category string `json:"category"`
	Brand    string `json:"brand"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	IsActive *bool  `json:"is_active"`
	IsDigital *bool `json:"is_digital"`
	IsFeatured *bool `json:"is_featured"`
	IsOnSale  *bool `json:"is_on_sale"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

// ========== CATEGORY DTOs ==========

// CreateCategoryRequest represents the request to create a new category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug" binding:"required"`
	ParentID    *uint  `json:"parent_id"`
	Image       string `json:"image"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	ParentID    *uint   `json:"parent_id"`
	Image       *string `json:"image"`
	Icon        *string `json:"icon"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

// CategoryResponse represents the category response
type CategoryResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Slug        string              `json:"slug"`
	ParentID    *uint               `json:"parent_id"`
	Parent      *CategoryResponse   `json:"parent,omitempty"`
	Children    []CategoryResponse  `json:"children,omitempty"`
	Image       string              `json:"image"`
	Icon        string              `json:"icon"`
	SortOrder   int                 `json:"sort_order"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// ListCategoriesResponse represents the paginated list of categories
type ListCategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total"`
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"`
}

// ========== VARIANT DTOs ==========

// CreateVariantRequest represents the request to create a new product variant
type CreateVariantRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	SKU       string  `json:"sku" binding:"required"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock" binding:"min=0"`
	Weight    float64 `json:"weight"`
	Color     string  `json:"color"`
	Size      string  `json:"size"`
	Material  string  `json:"material"`
	Image     string  `json:"image"`
	SortOrder int     `json:"sort_order"`
}

// UpdateVariantRequest represents the request to update a product variant
type UpdateVariantRequest struct {
	Name      *string  `json:"name"`
	SKU       *string  `json:"sku"`
	Price     *float64 `json:"price"`
	Stock     *int     `json:"stock"`
	Weight    *float64 `json:"weight"`
	Color     *string  `json:"color"`
	Size      *string  `json:"size"`
	Material  *string  `json:"material"`
	Image     *string  `json:"image"`
	IsActive  *bool    `json:"is_active"`
	SortOrder *int     `json:"sort_order"`
}

// VariantResponse represents the product variant response
type VariantResponse struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Weight    float64   `json:"weight"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
	Material  string    `json:"material"`
	Image     string    `json:"image"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListVariantsResponse represents the paginated list of variants
type ListVariantsResponse struct {
	Variants []VariantResponse `json:"variants"`
	Total    int               `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}
