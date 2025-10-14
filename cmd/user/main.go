package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/domain"
	"github.com/ddd-micro/internal/user/infrastructure/config"
	"github.com/ddd-micro/internal/user/infrastructure/database"
	"github.com/ddd-micro/internal/user/infrastructure/persistence"
	userhttp "github.com/ddd-micro/internal/user/interfaces/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := database.NewPostgresConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auto migrate database schema
	if err := db.GetDB().AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize repository
	userRepo := persistence.NewUserRepository(db.GetDB())

	// Initialize service
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	tokenDuration := 24 * time.Hour
	userService := application.NewUserService(userRepo, jwtSecret, tokenDuration)

	// Initialize Gin router
	ginMode := getEnv("GIN_MODE", "debug")
	gin.SetMode(ginMode)
	router := gin.Default()

	// Add CORS middleware
	router.Use(userhttp.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "user-service",
			"time":    time.Now().UTC(),
		})
	})

	// Setup routes
	userhttp.SetupRoutes(router, userService)

	// Server configuration
	port := getEnv("PORT", "8080")
	serverAddr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting User Service on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

