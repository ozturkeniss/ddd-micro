package application

import (
	"github.com/ddd-micro/internal/payment/application/command"
	"github.com/ddd-micro/internal/payment/application/query"
	"github.com/ddd-micro/internal/payment/domain"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for application layer
var ProviderSet = wire.NewSet(
	NewPaymentServiceCQRS,
	// Command handlers
	command.NewCreatePaymentCommandHandler,
	command.NewProcessPaymentCommandHandler,
	command.NewCancelPaymentCommandHandler,
	command.NewAddPaymentMethodCommandHandler,
	command.NewUpdatePaymentMethodCommandHandler,
	command.NewDeletePaymentMethodCommandHandler,
	// Query handlers
	query.NewGetPaymentQueryHandler,
	query.NewListPaymentsQueryHandler,
	query.NewGetPaymentMethodQueryHandler,
	query.NewListPaymentMethodsQueryHandler,
)
