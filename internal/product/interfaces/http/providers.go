package http

import (
	"github.com/ddd-micro/internal/product/application"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for HTTP interface layer
var ProviderSet = wire.NewSet(
	NewProductHandler,
	NewUserHandler,
	NewHTTPRouter,
)

// NewHTTPRouter creates a new HTTP router with all routes
func NewHTTPRouter(productHandler *ProductHandler, userHandler *UserHandler) *gin.Engine {
	router := gin.Default()

	// Setup routes
	SetupRoutes(router, productHandler, userHandler)

	return router
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *application.ProductServiceCQRS) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
