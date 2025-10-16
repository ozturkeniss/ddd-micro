package client

import (
	"github.com/ddd-micro/internal/basket/infrastructure/config"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for client layer
var ProviderSet = wire.NewSet(
	NewUserClientFromConfig,
)

// NewUserClientFromConfig creates a new user client from configuration
func NewUserClientFromConfig(cfg config.ClientConfig) UserClient {
	client, err := NewUserClient(cfg.UserService.URL)
	if err != nil {
		// In a real application, you might want to handle this error differently
		// For now, we'll panic since this is a critical dependency
		panic(err)
	}
	return client
}
