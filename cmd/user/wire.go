//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ddd-micro/internal/user/application"
	"github.com/ddd-micro/internal/user/infrastructure"
	"github.com/ddd-micro/internal/user/infrastructure/database"
	userhttp "github.com/ddd-micro/internal/user/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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

		// App constructor
		NewApp,
	)
	return &App{}, nil
}

// App holds all application dependencies
type App struct {
	Router      *gin.Engine
	UserService *application.UserServiceCQRS
	Database    *database.Database
}

// NewApp creates a new App instance
func NewApp(
	router *gin.Engine,
	userService *application.UserServiceCQRS,
	db *database.Database,
) *App {
	return &App{
		Router:      router,
		UserService: userService,
		Database:    db,
	}
}

