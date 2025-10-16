package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for HTTP interface layer
var ProviderSet = wire.NewSet(
	NewProductHandler,
	NewUserHandler,
	NewAuthMiddleware,
	NewHTTPRouter,
)

// NewHTTPRouter creates a new HTTP router with all routes
func NewHTTPRouter(productHandler *ProductHandler, userHandler *UserHandler, authMiddleware *AuthMiddleware) *gin.Engine {
	router := gin.Default()

	// Setup routes
	SetupRoutes(router, productHandler, userHandler, authMiddleware)

	return router
}

