package infrastructure

import (
	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/ddd-micro/internal/payment/infrastructure/config"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for infrastructure
var ProviderSet = wire.NewSet(
	config.LoadConfig,
	client.NewUserClient,
)
