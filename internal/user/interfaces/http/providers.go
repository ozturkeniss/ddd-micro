package http

import (
	"github.com/ddd-micro/internal/user/application"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for HTTP interface layer
var ProviderSet = wire.NewSet(
	ProvideUserHandler,
	ProvideRouter,
)

// ProvideUserHandler provides user HTTP handler
func ProvideUserHandler(userService *application.UserServiceCQRS) *UserHandler {
	return NewUserHandler(userService)
}

// ProvideRouter provides configured Gin router
func ProvideRouter(userService *application.UserServiceCQRS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // Can be overridden by environment
	router := gin.Default()

	// Add CORS middleware
	router.Use(CORSMiddleware())

	// Setup routes
	SetupRoutes(router, userService)

	return router
}

