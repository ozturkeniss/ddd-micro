package http

import (
	"github.com/ddd-micro/internal/payment/application"
	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Providers contains all HTTP-related dependencies
type Providers struct {
	Router *gin.Engine
}

// ProviderSet is the Wire provider set for HTTP
var ProviderSet = wire.NewSet(
	NewRouter,
	NewProviders,
)

// NewProviders creates new HTTP providers
func NewProviders(router *gin.Engine) *Providers {
	return &Providers{
		Router: router,
	}
}

// NewRouter creates a new Gin router with all routes configured
func NewRouter(
	paymentService *application.PaymentServiceCQRS,
	userClient client.UserClient,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Setup routes
	SetupRoutes(router, paymentService, userClient)

	return router
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
