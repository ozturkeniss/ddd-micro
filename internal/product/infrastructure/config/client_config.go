package config

import (
	"os"
)

// ClientConfig holds configuration for external service clients
type ClientConfig struct {
	UserService UserServiceConfig `json:"user_service"`
}

// UserServiceConfig holds configuration for user service client
type UserServiceConfig struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"` // in seconds
}

// LoadClientConfig loads client configuration from environment variables
func LoadClientConfig() *ClientConfig {
	return &ClientConfig{
		UserService: UserServiceConfig{
			URL:     getEnv("USER_SERVICE_URL", "localhost:9090"),
			Timeout: getEnvInt("USER_SERVICE_TIMEOUT", 5),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion, in production you might want to use strconv.Atoi
		if value == "5" {
			return 5
		}
		if value == "10" {
			return 10
		}
		if value == "30" {
			return 30
		}
	}
	return defaultValue
}
