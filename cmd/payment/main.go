// @title Payment Service API
// @version 1.0
// @description Payment Service API for managing payments, payment methods, and refunds
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8084
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ddd-micro/cmd/payment/docs" // This is required for swagger docs
	"github.com/ddd-micro/internal/payment/infrastructure/monitoring"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialize application
	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()
	defer app.JaegerTracer.Close()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Start HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	log.Printf("Starting HTTP Server on port %s...", httpPort)
	if err := app.HTTPRouter.Run(":" + httpPort); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	log.Println("Stopping HTTP server...")
	// Note: Gin doesn't have built-in graceful shutdown, but we can add it if needed

	log.Println("Server stopped")
}

// App represents the application dependencies
type App struct {
	HTTPRouter   *gin.Engine
	JaegerTracer *monitoring.JaegerTracer
}

// InitializeApp initializes all application dependencies using Wire
func InitializeApp() (*App, func(), error) {
	// For now, create a basic router setup
	router := gin.Default()

	// Add Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "payment-service"})
	})

	// Add basic API routes for Swagger generation
	api := router.Group("/api/v1")
	{
		// Payment routes
		payments := api.Group("/payments")
		{
			payments.POST("", func(c *gin.Context) {
				c.JSON(201, gin.H{"message": "Payment created"})
			})
			payments.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"payments": []gin.H{}})
			})
			payments.GET("/:id", func(c *gin.Context) {
				c.JSON(200, gin.H{"id": c.Param("id")})
			})
			payments.POST("/:id/process", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Payment processed"})
			})
			payments.POST("/:id/cancel", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Payment cancelled"})
			})
		}

		// Payment method routes
		paymentMethods := api.Group("/payment-methods")
		{
			paymentMethods.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"payment_methods": []gin.H{}})
			})
			paymentMethods.POST("", func(c *gin.Context) {
				c.JSON(201, gin.H{"message": "Payment method added"})
			})
			paymentMethods.PUT("/:id", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Payment method updated"})
			})
			paymentMethods.DELETE("/:id", func(c *gin.Context) {
				c.Status(204)
			})
			paymentMethods.POST("/:id/set-default", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Payment method set as default"})
			})
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			adminPayments := admin.Group("/payments")
			{
				adminPayments.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"payments": []gin.H{}})
				})
				adminPayments.GET("/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"id": c.Param("id")})
				})
				adminPayments.PUT("/:id/status", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Payment status updated"})
				})
			}

			adminRefunds := admin.Group("/refunds")
			{
				adminRefunds.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"refunds": []gin.H{}})
				})
				adminRefunds.POST("", func(c *gin.Context) {
					c.JSON(201, gin.H{"message": "Refund created"})
				})
				adminRefunds.GET("/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"id": c.Param("id")})
				})
				adminRefunds.POST("/:id/process", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Refund processed"})
				})
			}

			adminAnalytics := admin.Group("/analytics")
			{
				adminAnalytics.GET("/payments", func(c *gin.Context) {
					c.JSON(200, gin.H{"stats": gin.H{}})
				})
			}
		}
	}

	return &App{
			HTTPRouter: router,
		}, func() {
			// Cleanup function
		}, nil
}
