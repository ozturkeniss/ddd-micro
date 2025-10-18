package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application with Wire dependency injection
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Database.Close()
	defer app.JaegerTracer.Close()

	// Health check endpoint
	app.HTTPRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "user-service",
			"time":    time.Now().UTC(),
		})
	})

	// Server configuration
	httpPort := getEnv("HTTP_PORT", "8080")
	grpcPort := getEnv("GRPC_PORT", "9090")
	ginMode := getEnv("GIN_MODE", "debug")

	// Set Gin mode
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// HTTP Server
	httpAddr := fmt.Sprintf(":%s", httpPort)
	httpSrv := &http.Server{
		Addr:         httpAddr,
		Handler:      app.HTTPRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP Server on port %s...", httpPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// gRPC Server
	grpcAddr := fmt.Sprintf(":%s", grpcPort)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC Server on port %s...", grpcPort)
		if err := app.GRPCServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	log.Println("User Service started successfully!")
	log.Printf("  - HTTP API: http://localhost:%s", httpPort)
	log.Printf("  - gRPC API: localhost:%s", grpcPort)

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}

	// Gracefully stop gRPC server
	app.GRPCServer.GracefulStop()

	log.Println("Servers exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
