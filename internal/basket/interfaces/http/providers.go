package http

import (
	"github.com/ddd-micro/internal/basket/application"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/ddd-micro/internal/basket/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
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
	metrics *monitoring.PrometheusMetrics,
	tracer *monitoring.JaegerTracer,
) *gin.Engine {
	// Create router
	router := gin.Default()

	// Create handlers
	basketHandler := NewBasketHandler(basketService, metrics)
	userHandler := NewUserHandler(userClient)
	authMiddleware := NewAuthMiddleware(userClient)

	// Setup routes
	SetupRoutes(router, basketHandler, userHandler, authMiddleware, metrics, tracer)

	return router
}
