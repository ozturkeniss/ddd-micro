package infrastructure

import (
	"github.com/ddd-micro/internal/payment/infrastructure/client"
	"github.com/ddd-micro/internal/payment/infrastructure/config"
	"github.com/ddd-micro/internal/payment/infrastructure/database"
	"github.com/ddd-micro/internal/payment/infrastructure/gateway"
	"github.com/ddd-micro/internal/payment/infrastructure/persistence"
	"github.com/ddd-micro/kafka"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for infrastructure
var ProviderSet = wire.NewSet(
	// Configuration
	config.LoadConfig,

	// Database
	database.NewPostgresDB,

	// Repositories
	persistence.NewPaymentRepository,
	persistence.NewPaymentMethodRepository,
	persistence.NewRefundRepository,

	// External service clients
	client.NewUserClient,
	client.NewProductClient,
	client.NewBasketClient,

	// Payment gateways
	gateway.NewStripeGateway,
	gateway.NewMockGateway,

	// Kafka
	kafka.LoadConfig,
	kafka.NewKafkaPublisher,
)
