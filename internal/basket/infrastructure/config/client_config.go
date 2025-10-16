package config

// ClientConfig holds configuration for external service clients
type ClientConfig struct {
	UserService    UserServiceConfig
	ProductService ProductServiceConfig
}

// UserServiceConfig holds configuration for user service client
type UserServiceConfig struct {
	URL string
}

// ProductServiceConfig holds configuration for product service client
type ProductServiceConfig struct {
	URL string
}

// LoadClientConfig loads client configuration from environment variables
func LoadClientConfig() ClientConfig {
	return ClientConfig{
		UserService: UserServiceConfig{
			URL: getEnv("USER_SERVICE_URL", "localhost:9090"),
		},
		ProductService: ProductServiceConfig{
			URL: getEnv("PRODUCT_SERVICE_URL", "localhost:9091"),
		},
	}
}
