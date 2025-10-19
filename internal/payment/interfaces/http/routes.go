package http

import (
	"github.com/ddd-micro/internal/payment/application"
	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/ddd-micro/internal/payment/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes sets up all HTTP routes
func SetupRoutes(
	router *gin.Engine,
	paymentService *application.PaymentServiceCQRS,
	userClient client.UserClient,
	metrics *monitoring.PrometheusMetrics,
	tracer *monitoring.JaegerTracer,
) {
	// Add monitoring middlewares
	router.Use(monitoring.PrometheusMiddleware(metrics))
	router.Use(monitoring.JaegerMiddleware(tracer))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize handlers
	paymentHandler := NewPaymentHandler(paymentService, metrics)
	adminHandler := NewAdminHandler(paymentService)

	// Initialize middleware
	authMiddleware := AuthMiddleware(userClient)
	adminOnlyMiddleware := AdminOnlyMiddleware()
	userOrAdminMiddleware := UserOrAdminMiddleware()

	// Public routes (no authentication required)
	public := router.Group("/api/v1")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "payment-service"})
		})
	}

	// User routes (authentication required)
	user := router.Group("/api/v1")
	user.Use(authMiddleware)
	{
		// Payment routes
		payments := user.Group("/payments")
		{
			payments.POST("", paymentHandler.CreatePayment)              // POST /api/v1/payments
			payments.GET("", paymentHandler.ListPayments)                // GET /api/v1/payments
			payments.GET("/:id", paymentHandler.GetPayment)              // GET /api/v1/payments/:id
			payments.POST("/:id/process", paymentHandler.ProcessPayment) // POST /api/v1/payments/:id/process
			payments.POST("/:id/cancel", paymentHandler.CancelPayment)   // POST /api/v1/payments/:id/cancel
		}

		// Payment method routes
		paymentMethods := user.Group("/payment-methods")
		{
			paymentMethods.GET("", paymentHandler.GetPaymentMethods)                        // GET /api/v1/payment-methods
			paymentMethods.POST("", paymentHandler.AddPaymentMethod)                        // POST /api/v1/payment-methods
			paymentMethods.PUT("/:id", paymentHandler.UpdatePaymentMethod)                  // PUT /api/v1/payment-methods/:id
			paymentMethods.DELETE("/:id", paymentHandler.DeletePaymentMethod)               // DELETE /api/v1/payment-methods/:id
			paymentMethods.POST("/:id/set-default", paymentHandler.SetDefaultPaymentMethod) // POST /api/v1/payment-methods/:id/set-default
		}
	}

	// Admin routes (admin authentication required)
	admin := router.Group("/api/v1/admin")
	admin.Use(authMiddleware)
	admin.Use(adminOnlyMiddleware)
	{
		// Admin payment routes
		adminPayments := admin.Group("/payments")
		{
			adminPayments.GET("", adminHandler.ListAllPayments)                // GET /api/v1/admin/payments
			adminPayments.GET("/:id", adminHandler.GetPaymentByID)             // GET /api/v1/admin/payments/:id
			adminPayments.PUT("/:id/status", adminHandler.UpdatePaymentStatus) // PUT /api/v1/admin/payments/:id/status
		}

		// Admin refund routes
		adminRefunds := admin.Group("/refunds")
		{
			adminRefunds.GET("", adminHandler.ListRefunds)                // GET /api/v1/admin/refunds
			adminRefunds.POST("", adminHandler.CreateRefund)              // POST /api/v1/admin/refunds
			adminRefunds.GET("/:id", adminHandler.GetRefundByID)          // GET /api/v1/admin/refunds/:id
			adminRefunds.POST("/:id/process", adminHandler.ProcessRefund) // POST /api/v1/admin/refunds/:id/process
		}

		// Admin analytics routes
		adminAnalytics := admin.Group("/analytics")
		{
			adminAnalytics.GET("/payments", adminHandler.GetPaymentStats) // GET /api/v1/admin/analytics/payments
		}
	}

}
