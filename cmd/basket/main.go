package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ddd-micro/cmd/basket/docs"
)

// @title Basket Service API
// @version 1.0
// @description Basket Service API for managing user baskets
// @host localhost:8083
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	loadEnv()

	// Initialize application dependencies
	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer cleanup()

	// Setup graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP Server on port %s...", os.Getenv("HTTP_PORT"))
		if err := app.HTTPRouter.Run(":" + os.Getenv("HTTP_PORT")); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		grpcPort := os.Getenv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "9093"
		}
		
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
		}

		log.Printf("Starting gRPC Server on port %s...", grpcPort)
		if err := app.GRPCServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Gracefully stop gRPC server
	go func() {
		log.Println("Stopping gRPC server...")
		app.GRPCServer.GracefulStop()
	}()

	// Cancel main context
	cancel()

	// Wait for shutdown context
	<-shutdownCtx.Done()

	log.Println("Servers exited gracefully")
}

// loadEnv loads environment variables with defaults
func loadEnv() {
	if os.Getenv("HTTP_PORT") == "" {
		os.Setenv("HTTP_PORT", "8083")
	}
	if os.Getenv("GRPC_PORT") == "" {
		os.Setenv("GRPC_PORT", "9093")
	}
	if os.Getenv("REDIS_HOST") == "" {
		os.Setenv("REDIS_HOST", "localhost")
	}
	if os.Getenv("REDIS_PORT") == "" {
		os.Setenv("REDIS_PORT", "6379")
	}
	if os.Getenv("REDIS_PASSWORD") == "" {
		os.Setenv("REDIS_PASSWORD", "")
	}
	if os.Getenv("REDIS_DB") == "" {
		os.Setenv("REDIS_DB", "0")
	}
	if os.Getenv("USER_SERVICE_URL") == "" {
		os.Setenv("USER_SERVICE_URL", "localhost:9090")
	}
	if os.Getenv("PRODUCT_SERVICE_URL") == "" {
		os.Setenv("PRODUCT_SERVICE_URL", "localhost:9091")
	}
}
