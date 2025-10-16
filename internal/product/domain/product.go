package domain

import (
	"time"

	"gorm.io/gorm"
)

// Product represents the product domain entity
type Product struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"not null;size:255" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	ShortDescription string        `gorm:"size:500" json:"short_description"`
	Price           float64        `gorm:"not null;type:decimal(10,2)" json:"price"`
	ComparePrice    float64        `gorm:"type:decimal(10,2)" json:"compare_price"` // Original price for discount display
	CostPrice       float64        `gorm:"type:decimal(10,2)" json:"cost_price"`    // Cost price for profit calculation
	Stock           int            `gorm:"not null;default:0" json:"stock"`
	MinStock        int            `gorm:"default:0" json:"min_stock"` // Minimum stock alert
	MaxStock        int            `gorm:"default:0" json:"max_stock"` // Maximum stock limit
	Category        string         `gorm:"size:100" json:"category"`
	SubCategory     string         `gorm:"size:100" json:"sub_category"`
	Brand           string         `gorm:"size:100" json:"brand"`
	SKU             string         `gorm:"uniqueIndex;not null;size:100" json:"sku"`
	Barcode         string         `gorm:"size:50" json:"barcode"`
	Weight          float64        `gorm:"type:decimal(8,3)" json:"weight"` // Weight in kg
	Dimensions      string         `gorm:"size:100" json:"dimensions"`      // LxWxH format
	Color           string         `gorm:"size:50" json:"color"`
	Size            string         `gorm:"size:50" json:"size"`
	Material        string         `gorm:"size:100" json:"material"`
	Tags            string         `gorm:"type:text" json:"tags"` // Comma-separated tags
	Images          string         `gorm:"type:text" json:"images"` // JSON array of image URLs
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsDigital       bool           `gorm:"default:false" json:"is_digital"` // Digital product flag
	IsFeatured      bool           `gorm:"default:false" json:"is_featured"` // Featured product flag
	IsOnSale        bool           `gorm:"default:false" json:"is_on_sale"` // On sale flag
	SortOrder       int            `gorm:"default:0" json:"sort_order"` // For custom sorting
	ViewCount       int            `gorm:"default:0" json:"view_count"` // Product view counter
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
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

// IsValidBarcode checks if the barcode is valid
func (p *Product) IsValidBarcode() bool {
	return len(p.Barcode) == 0 || len(p.Barcode) <= 50
}

// IsValidWeight checks if the weight is valid
func (p *Product) IsValidWeight() bool {
	return p.Weight >= 0
}

// IsValidDimensions checks if the dimensions are valid
func (p *Product) IsValidDimensions() bool {
	return len(p.Dimensions) == 0 || len(p.Dimensions) <= 100
}

// IsValidShortDescription checks if the short description is valid
func (p *Product) IsValidShortDescription() bool {
	return len(p.ShortDescription) <= 500
}

// IsValidComparePrice checks if the compare price is valid
func (p *Product) IsValidComparePrice() bool {
	return p.ComparePrice >= 0 && p.ComparePrice >= p.Price
}

// IsValidCostPrice checks if the cost price is valid
func (p *Product) IsValidCostPrice() bool {
	return p.CostPrice >= 0
}

// IsValidMinStock checks if the minimum stock is valid
func (p *Product) IsValidMinStock() bool {
	return p.MinStock >= 0
}

// IsValidMaxStock checks if the maximum stock is valid
func (p *Product) IsValidMaxStock() bool {
	return p.MaxStock == 0 || p.MaxStock >= p.MinStock
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

// IsLowStock checks if the product stock is below minimum threshold
func (p *Product) IsLowStock() bool {
	return p.Stock <= p.MinStock && p.MinStock > 0
}

// IsOverStock checks if the product stock exceeds maximum limit
func (p *Product) IsOverStock() bool {
	return p.MaxStock > 0 && p.Stock > p.MaxStock
}

// GetDiscountPercentage calculates the discount percentage
func (p *Product) GetDiscountPercentage() float64 {
	if p.ComparePrice <= 0 || p.ComparePrice <= p.Price {
		return 0
	}
	return ((p.ComparePrice - p.Price) / p.ComparePrice) * 100
}

// GetProfitMargin calculates the profit margin
func (p *Product) GetProfitMargin() float64 {
	if p.CostPrice <= 0 {
		return 0
	}
	return ((p.Price - p.CostPrice) / p.CostPrice) * 100
}

// GetProfitAmount calculates the profit amount per unit
func (p *Product) GetProfitAmount() float64 {
	return p.Price - p.CostPrice
}

// MarkAsFeatured marks the product as featured
func (p *Product) MarkAsFeatured() {
	p.IsFeatured = true
}

// UnmarkAsFeatured removes featured status from the product
func (p *Product) UnmarkAsFeatured() {
	p.IsFeatured = false
}

// MarkAsOnSale marks the product as on sale
func (p *Product) MarkAsOnSale() {
	p.IsOnSale = true
}

// UnmarkAsOnSale removes on sale status from the product
func (p *Product) UnmarkAsOnSale() {
	p.IsOnSale = false
}

// IncrementViewCount increments the view count
func (p *Product) IncrementViewCount() {
	p.ViewCount++
}

// SetSortOrder sets the sort order for the product
func (p *Product) SetSortOrder(order int) {
	p.SortOrder = order
}

// IsDigitalProduct checks if the product is digital
func (p *Product) IsDigitalProduct() bool {
	return p.IsDigital
}

// IsPhysicalProduct checks if the product is physical
func (p *Product) IsPhysicalProduct() bool {
	return !p.IsDigital
}

// ValidateProduct validates all product fields
func (p *Product) ValidateProduct() error {
	if !p.IsValidName() {
		return ErrInvalidProductData
	}
	if !p.IsValidPrice() {
		return ErrInvalidProductData
	}
	if !p.IsValidStock() {
		return ErrInvalidProductData
	}
	if !p.IsValidSKU() {
		return ErrInvalidProductData
	}
	if !p.IsValidBarcode() {
		return ErrInvalidProductData
	}
	if !p.IsValidWeight() {
		return ErrInvalidProductData
	}
	if !p.IsValidDimensions() {
		return ErrInvalidProductData
	}
	if !p.IsValidShortDescription() {
		return ErrInvalidProductData
	}
	if !p.IsValidComparePrice() {
		return ErrInvalidProductData
	}
	if !p.IsValidCostPrice() {
		return ErrInvalidProductData
	}
	if !p.IsValidMinStock() {
		return ErrInvalidProductData
	}
	if !p.IsValidMaxStock() {
		return ErrInvalidProductData
	}
	return nil
}
