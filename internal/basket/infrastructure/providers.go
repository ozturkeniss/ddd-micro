package infrastructure

import (
	"log"

	"github.com/ddd-micro/internal/basket/domain"
	"github.com/ddd-micro/internal/basket/infrastructure/client"
	"github.com/ddd-micro/internal/basket/infrastructure/config"
	"github.com/ddd-micro/internal/basket/infrastructure/database"
	"github.com/ddd-micro/internal/basket/infrastructure/persistence"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for infrastructure layer
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideDatabaseConfig,
	ProvideDatabase,
	ProvideBasketRepository,
	ProvideUserClient,
	ProvideProductClient,
	NewConfig,
	NewRedisClient,
	NewUserClient,
	NewProductClient,
	NewBasketRepository,
)

// ProvideConfig provides application configuration
func ProvideConfig() *config.Config {
	return config.LoadConfig()
}

// ProvideDatabaseConfig provides database configuration
func ProvideDatabaseConfig(cfg *config.Config) database.Config {
	return database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Password: cfg.Database.Password,
		DB:       cfg.Database.DB,
	}
}

// ProvideDatabase provides Redis database connection
func ProvideDatabase(cfg database.Config) (*database.Database, error) {
	db, err := database.NewRedisConnection(cfg)
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Println("Redis connection established successfully")
	return db, nil
}

// ProvideBasketRepository provides basket repository
func ProvideBasketRepository(db *database.Database) domain.BasketRepository {
	return persistence.NewBasketRepository(db.GetClient())
}

// ProvideUserClient provides user client
func ProvideUserClient(cfg *config.Config) client.UserClient {
	return client.NewUserClientFromConfig(cfg.Client)
}

// ProvideProductClient provides product client
func ProvideProductClient(cfg *config.Config) client.ProductClient {
	return client.NewProductClientFromConfig(cfg.Client)
}

// NewConfig creates a new config
func NewConfig() *config.Config {
	return config.LoadConfig()
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config) *database.Database {
	db, err := database.NewRedisConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Password: cfg.Database.Password,
		DB:       cfg.Database.DB,
	})
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil
	}
	return db
}

// NewUserClient creates a new user client
func NewUserClient(cfg *config.Config) client.UserClient {
	return client.NewUserClientFromConfig(cfg.Client)
}

// NewProductClient creates a new product client
func NewProductClient(cfg *config.Config) client.ProductClient {
	return client.NewProductClientFromConfig(cfg.Client)
}

// NewBasketRepository creates a new basket repository
func NewBasketRepository(db *database.Database) domain.BasketRepository {
	return persistence.NewBasketRepository(db.GetClient())
}
