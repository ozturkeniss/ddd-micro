package domain

import (
	"time"
)

// Basket represents a shopping basket/cart
type Basket struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Items     []BasketItem `json:"items" gorm:"foreignKey:BasketID;constraint:OnDelete:CASCADE"`
	Total     float64   `json:"total" gorm:"type:decimal(10,2);default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index"`
}

// BasketItem represents an item in the basket
type BasketItem struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	BasketID  string  `json:"basket_id" gorm:"not null;index;type:varchar(36)"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Quantity  int     `json:"quantity" gorm:"not null;default:1"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice float64 `json:"total_price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name for Basket
func (Basket) TableName() string {
	return "baskets"
}

// TableName returns the table name for BasketItem
func (BasketItem) TableName() string {
	return "basket_items"
}

// CalculateTotal calculates the total price of the basket
func (b *Basket) CalculateTotal() {
	total := 0.0
	for _, item := range b.Items {
		total += item.TotalPrice
	}
	b.Total = total
}

// AddItem adds an item to the basket or updates quantity if exists
func (b *Basket) AddItem(productID uint, quantity int, unitPrice float64) {
	// Check if item already exists
	for i, item := range b.Items {
		if item.ProductID == productID {
			// Update existing item
			b.Items[i].Quantity += quantity
			b.Items[i].TotalPrice = float64(b.Items[i].Quantity) * b.Items[i].UnitPrice
			b.CalculateTotal()
			return
		}
	}
	
	// Add new item
	newItem := BasketItem{
		BasketID:   b.ID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: float64(quantity) * unitPrice,
	}
	b.Items = append(b.Items, newItem)
	b.CalculateTotal()
}

// RemoveItem removes an item from the basket
func (b *Basket) RemoveItem(productID uint) {
	for i, item := range b.Items {
		if item.ProductID == productID {
			// Remove item from slice
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
			b.CalculateTotal()
			return
		}
	}
}

// UpdateItemQuantity updates the quantity of an item
func (b *Basket) UpdateItemQuantity(productID uint, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	for i, item := range b.Items {
		if item.ProductID == productID {
			b.Items[i].Quantity = quantity
			b.Items[i].TotalPrice = float64(quantity) * b.Items[i].UnitPrice
			b.CalculateTotal()
			return nil
		}
	}
	
	return ErrItemNotFound
}

// Clear removes all items from the basket
func (b *Basket) Clear() {
	b.Items = []BasketItem{}
	b.Total = 0
}

// IsEmpty checks if the basket is empty
func (b *Basket) IsEmpty() bool {
	return len(b.Items) == 0
}

// GetItemCount returns the total number of items in the basket
func (b *Basket) GetItemCount() int {
	count := 0
	for _, item := range b.Items {
		count += item.Quantity
	}
	return count
}

// IsExpired checks if the basket has expired
func (b *Basket) IsExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// SetExpiration sets the expiration time for the basket
func (b *Basket) SetExpiration(duration time.Duration) {
	b.ExpiresAt = time.Now().Add(duration)
}

// GetItemByProductID finds an item by product ID
func (b *Basket) GetItemByProductID(productID uint) *BasketItem {
	for _, item := range b.Items {
		if item.ProductID == productID {
			return &item
		}
	}
	return nil
}

// Validate validates the basket
func (b *Basket) Validate() error {
	if b.UserID == 0 {
		return ErrInvalidUserID
	}
	
	if b.ID == "" {
		return ErrInvalidBasketID
	}
	
	for _, item := range b.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	
	return nil
}

// Validate validates the basket item
func (bi *BasketItem) Validate() error {
	if bi.ProductID == 0 {
		return ErrInvalidProductID
	}
	
	if bi.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	if bi.UnitPrice < 0 {
		return ErrInvalidPrice
	}
	
	if bi.BasketID == "" {
		return ErrInvalidBasketID
	}
	
	return nil
}
