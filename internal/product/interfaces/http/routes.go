package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes(router *gin.Engine, productHandler *ProductHandler, userHandler *UserHandler) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.ListProducts)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

		// User routes (for user validation)
		users := v1.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.POST("/validate-token", userHandler.ValidateToken)
		}
	}
}
