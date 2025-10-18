package http

import (
	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for HTTP interface layer
var ProviderSet = wire.NewSet(
	ProvideUserHandler,
	ProvideRouter,
)

// ProvideUserHandler provides user HTTP handler
func ProvideUserHandler(userService *application.UserServiceCQRS, metrics *monitoring.PrometheusMetrics) *UserHandler {
	return NewUserHandler(userService, metrics)
}

// ProvideRouter provides configured Gin router
func ProvideRouter(userService *application.UserServiceCQRS, metrics *monitoring.PrometheusMetrics, tracer *monitoring.JaegerTracer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Can be overridden by environment
	router := gin.Default()

	// Add CORS middleware
	router.Use(CORSMiddleware())

	// Setup routes
	SetupRoutes(router, userService, metrics, tracer)

	return router
}
