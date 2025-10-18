package main

import (
	"github.com/ddd-micro/internal/payment/application"
	"github.com/ddd-micro/internal/payment/infrastructure"
	"github.com/ddd-micro/internal/payment/infrastructure/kafka"
	"github.com/ddd-micro/internal/payment/interfaces/grpc"
	"github.com/ddd-micro/internal/payment/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// App represents the application dependencies
type App struct {
	HTTPRouter *gin.Engine
	GRPCServer *grpc.Server
}

// InitializeApp initializes all application dependencies using Wire
func InitializeApp() (*App, func(), error) {
	wire.Build(
		// Infrastructure layer
		infrastructure.ProviderSet,

		// Application layer
		application.ProviderSet,

		// HTTP interface layer
		http.ProviderSet,

		// gRPC interface layer
		grpc.ProviderSet,

		// Kafka layer
		kafka.ProviderSet,

		// Main app
		NewApp,
	)

	return &App{}, nil, nil
}

// NewApp creates a new App instance
func NewApp(httpRouter *gin.Engine, grpcServer *grpc.Server) *App {
	return &App{
		HTTPRouter: httpRouter,
		GRPCServer: grpcServer,
	}
}
