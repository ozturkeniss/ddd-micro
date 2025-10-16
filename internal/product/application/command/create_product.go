package command

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// CreateProductCommand represents the command to create a new product
type CreateProductCommand struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Price            float64 `json:"price"`
	ComparePrice     float64 `json:"compare_price"`
	CostPrice        float64 `json:"cost_price"`
	Stock            int     `json:"stock"`
	MinStock         int     `json:"min_stock"`
	MaxStock         int     `json:"max_stock"`
	Category         string  `json:"category"`
	SubCategory      string  `json:"sub_category"`
	Brand            string  `json:"brand"`
	SKU              string  `json:"sku"`
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

// CreateProductHandler handles the create product command
type CreateProductHandler struct {
	repo domain.ProductRepository
}

// NewCreateProductHandler creates a new create product handler
func NewCreateProductHandler(repo domain.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{
		repo: repo,
	}
}

// Handle executes the create product command
func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*domain.Product, error) {
	// Create product entity
	product := &domain.Product{
		Name:             cmd.Name,
		Description:      cmd.Description,
		ShortDescription: cmd.ShortDescription,
		Price:            cmd.Price,
		ComparePrice:     cmd.ComparePrice,
		CostPrice:        cmd.CostPrice,
		Stock:            cmd.Stock,
		MinStock:         cmd.MinStock,
		MaxStock:         cmd.MaxStock,
		Category:         cmd.Category,
		SubCategory:      cmd.SubCategory,
		Brand:            cmd.Brand,
		SKU:              cmd.SKU,
		Barcode:          cmd.Barcode,
		Weight:           cmd.Weight,
		Dimensions:       cmd.Dimensions,
		Color:            cmd.Color,
		Size:             cmd.Size,
		Material:         cmd.Material,
		Tags:             cmd.Tags,
		Images:           cmd.Images,
		IsDigital:        cmd.IsDigital,
		IsFeatured:       cmd.IsFeatured,
		IsOnSale:         cmd.IsOnSale,
		SortOrder:        cmd.SortOrder,
		IsActive:         true,
	}

	// Validate product
	if err := product.ValidateProduct(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := h.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}
