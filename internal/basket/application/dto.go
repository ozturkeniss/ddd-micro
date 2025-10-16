package application

import "time"

// CreateBasketRequest represents the request to create a basket
type CreateBasketRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// AddItemRequest represents the request to add an item to the basket
type AddItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
}

// UpdateItemRequest represents the request to update an item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// BasketResponse represents the response for basket operations
type BasketResponse struct {
	ID        string             `json:"id"`
	UserID    uint               `json:"user_id"`
	Items     []BasketItemResponse `json:"items"`
	Total     float64            `json:"total"`
	ItemCount int                `json:"item_count"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	ExpiresAt time.Time          `json:"expires_at"`
	IsExpired bool               `json:"is_expired"`
}

// BasketItemResponse represents the response for basket item operations
type BasketItemResponse struct {
	ID         uint    `json:"id"`
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ListBasketsResponse represents the response for listing baskets
type ListBasketsResponse struct {
	Baskets []BasketResponse `json:"baskets"`
	Total   int              `json:"total"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
