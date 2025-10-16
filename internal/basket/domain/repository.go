package domain

import "context"

// BasketRepository defines the interface for basket data operations
type BasketRepository interface {
	// Create creates a new basket
	Create(ctx context.Context, basket *Basket) error
	
	// GetByID retrieves a basket by ID
	GetByID(ctx context.Context, basketID string) (*Basket, error)
	
	// GetByUserID retrieves a basket by user ID
	GetByUserID(ctx context.Context, userID uint) (*Basket, error)
	
	// Update updates an existing basket
	Update(ctx context.Context, basket *Basket) error
	
	// Delete deletes a basket by ID
	Delete(ctx context.Context, basketID string) error
	
	// DeleteByUserID deletes a basket by user ID
	DeleteByUserID(ctx context.Context, userID uint) error
	
	// AddItem adds an item to the basket
	AddItem(ctx context.Context, basketID string, item *BasketItem) error
	
	// UpdateItem updates a basket item
	UpdateItem(ctx context.Context, basketID string, item *BasketItem) error
	
	// RemoveItem removes an item from the basket
	RemoveItem(ctx context.Context, basketID string, productID uint) error
	
	// ClearItems removes all items from the basket
	ClearItems(ctx context.Context, basketID string) error
	
	// Exists checks if a basket exists by ID
	Exists(ctx context.Context, basketID string) (bool, error)
	
	// ExistsByUserID checks if a basket exists for a user
	ExistsByUserID(ctx context.Context, userID uint) (bool, error)
	
	// CleanupExpired removes expired baskets
	CleanupExpired(ctx context.Context) error
	
	// GetExpiredBaskets retrieves expired baskets
	GetExpiredBaskets(ctx context.Context) ([]*Basket, error)
}
