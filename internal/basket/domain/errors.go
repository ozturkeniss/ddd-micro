package domain

import "errors"

var (
	// Basket errors
	ErrBasketNotFound     = errors.New("basket not found")
	ErrBasketExpired      = errors.New("basket has expired")
	ErrInvalidBasketID    = errors.New("invalid basket ID")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrBasketAlreadyExists = errors.New("basket already exists for this user")
	
	// BasketItem errors
	ErrItemNotFound       = errors.New("item not found in basket")
	ErrInvalidProductID   = errors.New("invalid product ID")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrInsufficientStock  = errors.New("insufficient stock")
	
	// General errors
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrServiceUnavailable = errors.New("service unavailable")
)
