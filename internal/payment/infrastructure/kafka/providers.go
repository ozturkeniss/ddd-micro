package kafka

import (
	"github.com/ddd-micro/kafka"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for Kafka
var ProviderSet = wire.NewSet(
	NewPaymentEventPublisher,
)

// NewPaymentEventPublisher creates a new payment event publisher
func NewPaymentEventPublisher(publisher kafka.EventPublisher) *PaymentEventPublisher {
	return &PaymentEventPublisher{
		publisher: publisher,
	}
}
