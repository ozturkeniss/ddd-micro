package http

import (
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

