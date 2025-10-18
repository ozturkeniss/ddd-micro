package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the payment service
type Config struct {
	// Server configuration
	HTTPPort string
	GRPCPort string

	// Database configuration
	Database DatabaseConfig

	// External services
	UserServiceURL    string
	ProductServiceURL string
	BasketServiceURL  string

	// Payment gateway configuration
	Stripe StripeConfig

	// JWT configuration
	JWT JWTConfig

	// Redis configuration (for caching)
	Redis RedisConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// StripeConfig holds Stripe payment gateway configuration
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8084"),
		GRPCPort: getEnv("GRPC_PORT", "9094"),

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "payment_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},

		UserServiceURL:    getEnv("USER_SERVICE_URL", "user-service:9091"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "product-service:9092"),
		BasketServiceURL:  getEnv("BASKET_SERVICE_URL", "basket-service:9093"),

		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},

		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key"),
			ExpiresIn: time.Duration(getEnvAsInt("JWT_EXPIRES_IN", 24)) * time.Hour,
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	return config, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
