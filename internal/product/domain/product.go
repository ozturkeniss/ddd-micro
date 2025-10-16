package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product represents the product domain entity
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Category    string         `gorm:"size:100" json:"category"`
	SKU         string         `gorm:"uniqueIndex;not null;size:100" json:"sku"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Product entity
func (Product) TableName() string {
	return "products"
}

// IsValidName checks if the product name is valid
func (p *Product) IsValidName() bool {
	return len(p.Name) > 0 && len(p.Name) <= 255
}

// IsValidPrice checks if the price is valid
func (p *Product) IsValidPrice() bool {
	return p.Price >= 0
}

// IsValidStock checks if the stock is valid
func (p *Product) IsValidStock() bool {
	return p.Stock >= 0
}

// IsValidSKU checks if the SKU is valid
func (p *Product) IsValidSKU() bool {
	return len(p.SKU) > 0 && len(p.SKU) <= 100
}

// Activate activates the product
func (p *Product) Activate() {
	p.IsActive = true
}

// Deactivate deactivates the product
func (p *Product) Deactivate() {
	p.IsActive = false
}

// ReduceStock reduces the stock by the specified amount
func (p *Product) ReduceStock(amount int) error {
	if amount <= 0 {
		return ErrInvalidStockAmount
	}
	if p.Stock < amount {
		return ErrInsufficientStock
	}
	p.Stock -= amount
	return nil
}

// IncreaseStock increases the stock by the specified amount
func (p *Product) IncreaseStock(amount int) error {
	if amount <= 0 {
		return ErrInvalidStockAmount
	}
	p.Stock += amount
	return nil
}

// IsInStock checks if the product is in stock
func (p *Product) IsInStock() bool {
	return p.Stock > 0 && p.IsActive
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.IsActive && p.IsInStock()
}
