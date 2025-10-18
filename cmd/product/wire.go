//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ddd-micro/internal/product/application"
	"github.com/ddd-micro/internal/product/infrastructure"
	"github.com/ddd-micro/internal/product/infrastructure/database"
	"github.com/ddd-micro/internal/product/infrastructure/monitoring"
	productgrpc "github.com/ddd-micro/internal/product/interfaces/grpc"
	producthttp "github.com/ddd-micro/internal/product/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// InitializeApp initializes the entire application with all dependencies
func InitializeApp() (*App, error) {
	wire.Build(
		// Infrastructure providers
		infrastructure.ProviderSet,

		// Application providers
		application.ProviderSet,

		// HTTP interface providers
		producthttp.ProviderSet,

		// gRPC interface providers
		productgrpc.ProviderSet,

		// App constructor
		NewApp,
	)
	return &App{}, nil
}

// App holds all application dependencies
type App struct {
	HTTPRouter     *gin.Engine
	GRPCServer     *grpc.Server
	ProductService *application.ProductServiceCQRS
	UserService    *application.UserService
	Database       *database.Database
	UserClient     interface{ Close() error }
	JaegerTracer   *monitoring.JaegerTracer
}

// NewApp creates a new App instance
func NewApp(
	httpRouter *gin.Engine,
	grpcServer *grpc.Server,
	productService *application.ProductServiceCQRS,
	userService *application.UserService,
	db *database.Database,
	userClient interface{ Close() error },
	jaegerTracer *monitoring.JaegerTracer,
) *App {
	return &App{
		HTTPRouter:     httpRouter,
		GRPCServer:     grpcServer,
		ProductService: productService,
		UserService:    userService,
		Database:       db,
		UserClient:     userClient,
		JaegerTracer:   jaegerTracer,
	}
}
