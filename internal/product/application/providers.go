package application

import (
	"github.com/google/wire"
)

// ProviderSet is the application layer providers
var ProviderSet = wire.NewSet(
	NewProductService,
	NewProductServiceCQRS,
	NewUserService,
)
