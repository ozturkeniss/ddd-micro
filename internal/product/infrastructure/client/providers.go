package client

import (
	"github.com/ddd-micro/internal/product/infrastructure/config"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for client layer
var ProviderSet = wire.NewSet(
	NewUserClient,
)

// NewUserClient creates a new user service client
func NewUserClient(config *config.ClientConfig) (UserClient, error) {
	return NewUserClient(config.UserService.URL)
}
