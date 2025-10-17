package http

import (
	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for HTTP interface layer
var ProviderSet = wire.NewSet(
	NewBasketHandler,
	NewUserHandler,
	NewAuthMiddleware,
	NewHTTPRouter,
)

// NewHTTPRouter creates and configures the HTTP router
func NewHTTPRouter(
	basketService *application.BasketServiceCQRS,
	userClient client.UserClient,
) *gin.Engine {
	// Create router
	router := gin.Default()

	// Create handlers
	basketHandler := NewBasketHandler(basketService)
	userHandler := NewUserHandler(userClient)
	authMiddleware := NewAuthMiddleware(userClient)

	// Setup routes
	SetupRoutes(router, basketHandler, userHandler, authMiddleware)

	return router
}
