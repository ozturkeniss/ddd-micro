package config

import "os"

// ClientConfig holds configuration for external service clients
type ClientConfig struct {
	UserService UserServiceConfig
}

// UserServiceConfig holds configuration for user service client
type UserServiceConfig struct {
	URL string
}

// LoadClientConfig loads client configuration from environment variables
func LoadClientConfig() ClientConfig {
	return ClientConfig{
		UserService: UserServiceConfig{
			URL: getEnv("USER_SERVICE_URL", "localhost:9090"),
		},
	}
}
