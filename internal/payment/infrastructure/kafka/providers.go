package kafka

import (
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for Kafka
var ProviderSet = wire.NewSet(
	NewPaymentEventPublisher,
)
