package domain

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product with this SKU already exists")
	ErrInvalidStockAmount   = errors.New("invalid stock amount")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidProductData   = errors.New("invalid product data")
	ErrProductNotActive     = errors.New("product is not active")
)
