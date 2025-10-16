package application

import (
	"github.com/ddd-micro/internal/basket/application/command"
	"github.com/ddd-micro/internal/basket/application/query"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for application layer
var ProviderSet = wire.NewSet(
	// Command handlers
	command.NewCreateBasketCommandHandler,
	command.NewAddItemCommandHandler,
	command.NewUpdateItemCommandHandler,
	command.NewRemoveItemCommandHandler,
	command.NewClearBasketCommandHandler,
	
	// Query handlers
	query.NewGetBasketQueryHandler,
	
	// Main service
	NewBasketServiceCQRS,
)