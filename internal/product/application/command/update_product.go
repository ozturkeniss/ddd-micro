package command

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// UpdateProductCommand represents the command to update a product
type UpdateProductCommand struct {
	ProductID        uint     `json:"product_id"`
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

// UpdateProductHandler handles the update product command
type UpdateProductHandler struct {
	repo domain.ProductRepository
}

// NewUpdateProductHandler creates a new update product handler
func NewUpdateProductHandler(repo domain.ProductRepository) *UpdateProductHandler {
	return &UpdateProductHandler{
		repo: repo,
	}
}

// Handle executes the update product command
func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) (*domain.Product, error) {
	// Get existing product
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if cmd.Name != nil {
		product.Name = *cmd.Name
	}
	if cmd.Description != nil {
		product.Description = *cmd.Description
	}
	if cmd.ShortDescription != nil {
		product.ShortDescription = *cmd.ShortDescription
	}
	if cmd.Price != nil {
		product.Price = *cmd.Price
	}
	if cmd.ComparePrice != nil {
		product.ComparePrice = *cmd.ComparePrice
	}
	if cmd.CostPrice != nil {
		product.CostPrice = *cmd.CostPrice
	}
	if cmd.Stock != nil {
		product.Stock = *cmd.Stock
	}
	if cmd.MinStock != nil {
		product.MinStock = *cmd.MinStock
	}
	if cmd.MaxStock != nil {
		product.MaxStock = *cmd.MaxStock
	}
	if cmd.Category != nil {
		product.Category = *cmd.Category
	}
	if cmd.SubCategory != nil {
		product.SubCategory = *cmd.SubCategory
	}
	if cmd.Brand != nil {
		product.Brand = *cmd.Brand
	}
	if cmd.Barcode != nil {
		product.Barcode = *cmd.Barcode
	}
	if cmd.Weight != nil {
		product.Weight = *cmd.Weight
	}
	if cmd.Dimensions != nil {
		product.Dimensions = *cmd.Dimensions
	}
	if cmd.Color != nil {
		product.Color = *cmd.Color
	}
	if cmd.Size != nil {
		product.Size = *cmd.Size
	}
	if cmd.Material != nil {
		product.Material = *cmd.Material
	}
	if cmd.Tags != nil {
		product.Tags = *cmd.Tags
	}
	if cmd.Images != nil {
		product.Images = *cmd.Images
	}
	if cmd.IsActive != nil {
		product.IsActive = *cmd.IsActive
	}
	if cmd.IsDigital != nil {
		product.IsDigital = *cmd.IsDigital
	}
	if cmd.IsFeatured != nil {
		product.IsFeatured = *cmd.IsFeatured
	}
	if cmd.IsOnSale != nil {
		product.IsOnSale = *cmd.IsOnSale
	}
	if cmd.SortOrder != nil {
		product.SortOrder = *cmd.SortOrder
	}

	// Validate updated product
	if err := product.ValidateProduct(); err != nil {
		return nil, err
	}

	// Save changes
	if err := h.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
