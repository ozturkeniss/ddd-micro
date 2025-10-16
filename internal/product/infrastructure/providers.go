package infrastructure

import (
	"github.com/ddd-micro/internal/product/infrastructure/client"
	"github.com/ddd-micro/internal/product/infrastructure/config"
	"github.com/ddd-micro/internal/product/infrastructure/database"
	"github.com/ddd-micro/internal/product/infrastructure/persistence"
	"github.com/google/wire"
)

// ProviderSet is the infrastructure layer providers
var ProviderSet = wire.NewSet(
	// Config providers
	config.LoadConfig,
	config.LoadClientConfig,

	// Database providers
	database.NewPostgresConnection,

	// Persistence providers
	persistence.NewProductRepository,

	// Client providers
	client.ProviderSet,
)
