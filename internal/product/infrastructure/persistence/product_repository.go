package persistence

import (
	"context"
	"errors"

	"github.com/ddd-micro/internal/product/domain"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product with this SKU already exists")
)

// ProductRepository is the concrete implementation of domain.ProductRepository
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create creates a new product in the database
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	// Check if product already exists
	exists, err := r.Exists(ctx, product.SKU)
	if err != nil {
		return err
	}
	if exists {
		return ErrProductAlreadyExists
	}

	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	result := r.db.WithContext(ctx).First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, result.Error
	}

	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product
	result := r.db.WithContext(ctx).Where("sku = ?", sku).First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, result.Error
	}

	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}

// Delete soft deletes a product by ID
func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}

// List retrieves all products with pagination
func (r *ProductRepository) List(ctx context.Context, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product

	result := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

// ListByCategory retrieves products by category with pagination
func (r *ProductRepository) ListByCategory(ctx context.Context, category string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product

	result := r.db.WithContext(ctx).
		Where("category = ?", category).
		Offset(offset).
		Limit(limit).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

// SearchByName searches products by name with pagination
func (r *ProductRepository) SearchByName(ctx context.Context, name string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product

	result := r.db.WithContext(ctx).
		Where("name ILIKE ?", "%"+name+"%").
		Offset(offset).
		Limit(limit).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

// Exists checks if a product exists by SKU
func (r *ProductRepository) Exists(ctx context.Context, sku string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("sku = ?", sku).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// UpdateStock updates the stock of a product
func (r *ProductRepository) UpdateStock(ctx context.Context, id uint, stock int) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", id).
		Update("stock", stock)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}
