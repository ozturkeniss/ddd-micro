package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all HTTP routes for the basket service
func SetupRoutes(router *gin.Engine, basketHandler *BasketHandler, userHandler *UserHandler, authMiddleware *AuthMiddleware) {
	// Health check endpoint (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "basket-service",
		})
	})

	// Swagger documentation endpoint (public)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// User routes (require authentication)
		users := v1.Group("/users")
		users.Use(authMiddleware.AuthRequired())
		{
			// User profile
			users.GET("/profile", userHandler.GetProfile)

			// Basket operations
			users.POST("/basket", basketHandler.CreateBasket)
			users.GET("/basket", basketHandler.GetBasket)
			users.POST("/basket/items", basketHandler.AddItem)
			users.PUT("/basket/items", basketHandler.UpdateItem)
			users.DELETE("/basket/items/:product_id", basketHandler.RemoveItem)
			users.DELETE("/basket/clear", basketHandler.ClearBasket)
		}

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.AdminRequired())
		{
			// Admin basket operations
			admin.GET("/baskets/:user_id", basketHandler.AdminGetBasket)
			admin.DELETE("/baskets/:user_id", basketHandler.AdminDeleteBasket)
			admin.POST("/baskets/cleanup", basketHandler.AdminCleanupExpiredBaskets)
		}

		// Public routes (no authentication required)
		public := v1.Group("/public")
		{
			// Token validation (public endpoint for token verification)
			public.POST("/validate-token", userHandler.ValidateToken)
		}
	}
}
