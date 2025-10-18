//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/infrastructure"
	"github.com/ddd-micro/internal/user/infrastructure/database"
	"github.com/ddd-micro/internal/user/infrastructure/monitoring"
	usergrpc "github.com/ddd-micro/internal/user/interfaces/grpc"
	userhttp "github.com/ddd-micro/internal/user/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// InitializeApp initializes the entire application with all dependencies
func InitializeApp() (*App, error) {
	wire.Build(
		// Infrastructure providers
		infrastructure.ProviderSet,
		infrastructure.ProvideJWTSecret,

		// Application providers
		application.ProviderSet,

		// HTTP interface providers
		userhttp.ProviderSet,

		// gRPC interface providers
		usergrpc.ProviderSet,

		// App constructor
		NewApp,
	)
	return &App{}, nil
}

// App holds all application dependencies
type App struct {
	HTTPRouter   *gin.Engine
	GRPCServer   *grpc.Server
	UserService  *application.UserServiceCQRS
	Database     *database.Database
	JaegerTracer *monitoring.JaegerTracer
}

// NewApp creates a new App instance
func NewApp(
	httpRouter *gin.Engine,
	grpcServer *grpc.Server,
	userService *application.UserServiceCQRS,
	db *database.Database,
	jaegerTracer *monitoring.JaegerTracer,
) *App {
	return &App{
		HTTPRouter:   httpRouter,
		GRPCServer:   grpcServer,
		UserService:  userService,
		Database:     db,
		JaegerTracer: jaegerTracer,
	}
}
