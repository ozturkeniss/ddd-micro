package infrastructure

import (
	"log"
	"os"

	"github.com/ddd-micro/internal/user/domain"
	"github.com/ddd-micro/internal/user/infrastructure/config"
	"github.com/ddd-micro/internal/user/infrastructure/database"
	"github.com/ddd-micro/internal/user/infrastructure/persistence"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for infrastructure layer
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideDatabaseConfig,
	ProvideDatabase,
	ProvideUserRepository,
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
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
}

// ProvideDatabase provides database connection
func ProvideDatabase(cfg database.Config) (*database.Database, error) {
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Auto migrate
	if err := db.GetDB().AutoMigrate(&domain.User{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}
	log.Println("Database migration completed successfully")

	return db, nil
}

// ProvideUserRepository provides user repository
func ProvideUserRepository(db *database.Database) domain.UserRepository {
	return persistence.NewUserRepository(db.GetDB())
}

// ProvideJWTSecret provides JWT secret key from environment
func ProvideJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return secret
}

