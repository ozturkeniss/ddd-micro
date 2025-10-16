package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProductVariant represents the product variant domain entity
type ProductVariant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductID   uint           `gorm:"not null;index" json:"product_id"`
	Product     *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	SKU         string         `gorm:"uniqueIndex;not null;size:100" json:"sku"`
	Price       float64        `gorm:"type:decimal(10,2)" json:"price"` // Override product price if set
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Weight      float64        `gorm:"type:decimal(8,3)" json:"weight"` // Override product weight if set
	Color       string         `gorm:"size:50" json:"color"`
	Size        string         `gorm:"size:50" json:"size"`
	Material    string         `gorm:"size:100" json:"material"`
	Image       string         `gorm:"size:500" json:"image"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ProductVariant entity
func (ProductVariant) TableName() string {
	return "product_variants"
}

// IsValidName checks if the variant name is valid
func (v *ProductVariant) IsValidName() bool {
	return len(v.Name) > 0 && len(v.Name) <= 255
}

// IsValidSKU checks if the SKU is valid
func (v *ProductVariant) IsValidSKU() bool {
	return len(v.SKU) > 0 && len(v.SKU) <= 100
}

// IsValidPrice checks if the price is valid
func (v *ProductVariant) IsValidPrice() bool {
	return v.Price >= 0
}

// IsValidStock checks if the stock is valid
func (v *ProductVariant) IsValidStock() bool {
	return v.Stock >= 0
}

// IsValidWeight checks if the weight is valid
func (v *ProductVariant) IsValidWeight() bool {
	return v.Weight >= 0
}

// GetEffectivePrice returns the variant price if set, otherwise product price
func (v *ProductVariant) GetEffectivePrice() float64 {
	if v.Price > 0 {
		return v.Price
	}
	if v.Product != nil {
		return v.Product.Price
	}
	return 0
}

// GetEffectiveWeight returns the variant weight if set, otherwise product weight
func (v *ProductVariant) GetEffectiveWeight() float64 {
	if v.Weight > 0 {
		return v.Weight
	}
	if v.Product != nil {
		return v.Product.Weight
	}
	return 0
}

// IsInStock checks if the variant is in stock
func (v *ProductVariant) IsInStock() bool {
	return v.Stock > 0 && v.IsActive
}

// IsAvailable checks if the variant is available for purchase
func (v *ProductVariant) IsAvailable() bool {
	return v.IsActive && v.IsInStock()
}

// ReduceStock reduces the stock by the specified amount
func (v *ProductVariant) ReduceStock(amount int) error {
	if amount <= 0 {
		return ErrInvalidStockAmount
	}
	if v.Stock < amount {
		return ErrInsufficientStock
	}
	v.Stock -= amount
	return nil
}

// IncreaseStock increases the stock by the specified amount
func (v *ProductVariant) IncreaseStock(amount int) error {
	if amount <= 0 {
		return ErrInvalidStockAmount
	}
	v.Stock += amount
	return nil
}

// Activate activates the variant
func (v *ProductVariant) Activate() {
	v.IsActive = true
}

// Deactivate deactivates the variant
func (v *ProductVariant) Deactivate() {
	v.IsActive = false
}

// SetSortOrder sets the sort order for the variant
func (v *ProductVariant) SetSortOrder(order int) {
	v.SortOrder = order
}

// ValidateVariant validates all variant fields
func (v *ProductVariant) ValidateVariant() error {
	if !v.IsValidName() {
		return ErrInvalidProductData
	}
	if !v.IsValidSKU() {
		return ErrInvalidProductData
	}
	if !v.IsValidPrice() {
		return ErrInvalidProductData
	}
	if !v.IsValidStock() {
		return ErrInvalidProductData
	}
	if !v.IsValidWeight() {
		return ErrInvalidProductData
	}
	return nil
}
