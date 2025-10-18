package http

import (
	"github.com/ddd-micro/internal/product/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRoutes sets up all HTTP routes with RBAC
func SetupRoutes(router *gin.Engine, productHandler *ProductHandler, userHandler *UserHandler, authMiddleware *AuthMiddleware, metrics *monitoring.PrometheusMetrics, tracer *monitoring.JaegerTracer) {
	// Add monitoring middlewares
	router.Use(monitoring.PrometheusMiddleware(metrics))
	router.Use(monitoring.JaegerMiddleware(tracer))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public product routes (no authentication required)
		public := v1.Group("/products")
		{
			public.GET("", productHandler.ListProducts)
			public.GET("/search", productHandler.SearchProducts)
			public.GET("/category/:category", productHandler.ListProductsByCategory)
			public.GET("/:id", productHandler.GetProduct)
			public.POST("/:id/view", productHandler.IncrementViewCount)
		}

		// User routes (authentication required)
		users := v1.Group("/users")
		users.Use(authMiddleware.AuthRequired())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.POST("/validate-token", userHandler.ValidateToken)
		}

		// Admin product routes (admin access required)
		admin := v1.Group("/admin/products")
		admin.Use(authMiddleware.AdminRequired())
		{
			admin.POST("", productHandler.CreateProduct)
			admin.PUT("/:id", productHandler.UpdateProduct)
			admin.DELETE("/:id", productHandler.DeleteProduct)
			admin.PUT("/:id/stock", productHandler.UpdateStock)
			admin.POST("/:id/reduce-stock", productHandler.ReduceStock)
			admin.POST("/:id/increase-stock", productHandler.IncreaseStock)
			admin.POST("/:id/activate", productHandler.ActivateProduct)
			admin.POST("/:id/deactivate", productHandler.DeactivateProduct)
			admin.POST("/:id/featured", productHandler.MarkAsFeatured)
			admin.DELETE("/:id/featured", productHandler.UnmarkAsFeatured)
		}
	}
}
