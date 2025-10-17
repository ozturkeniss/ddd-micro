package application

import (
	"github.com/ddd-micro/internal/payment/domain"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for application layer
var ProviderSet = wire.NewSet(
	NewPaymentServiceCQRS,
)
