package application

import "errors"

// Application layer errors
var (
	ErrBasketNotFound   = errors.New("basket not found")
	ErrItemNotFound     = errors.New("item not found in basket")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidProductID = errors.New("invalid product ID")
	ErrBasketExpired    = errors.New("basket has expired")
	ErrProductNotFound  = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)
