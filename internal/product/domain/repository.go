package domain

import "context"

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *Product) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id uint) (*Product, error)

	// GetBySKU retrieves a product by SKU
	GetBySKU(ctx context.Context, sku string) (*Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *Product) error

	// Delete soft deletes a product
	Delete(ctx context.Context, id uint) error

	// List retrieves all products with pagination
	List(ctx context.Context, offset, limit int) ([]*Product, error)

	// ListByCategory retrieves products by category with pagination
	ListByCategory(ctx context.Context, category string, offset, limit int) ([]*Product, error)

	// SearchByName searches products by name with pagination
	SearchByName(ctx context.Context, name string, offset, limit int) ([]*Product, error)

	// Exists checks if a product exists by SKU
	Exists(ctx context.Context, sku string) (bool, error)

	// UpdateStock updates the stock of a product
	UpdateStock(ctx context.Context, id uint, stock int) error
}
