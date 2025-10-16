package application

import (
	"github.com/ddd-micro/internal/product/application/command"
	"github.com/ddd-micro/internal/product/application/query"
	"github.com/ddd-micro/internal/product/domain"
	"github.com/google/wire"
)

// ProviderSet is the application layer providers
var ProviderSet = wire.NewSet(
	NewProductService,
	NewProductServiceCQRS,
)

// NewProductService creates a new product service (non-CQRS)
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// NewProductServiceCQRS creates a new CQRS-based product service
func NewProductServiceCQRS(repo domain.ProductRepository) *ProductServiceCQRS {
	return &ProductServiceCQRS{
		// Command handlers
		createProductHandler:        command.NewCreateProductHandler(repo),
		updateProductHandler:        command.NewUpdateProductHandler(repo),
		deleteProductHandler:        command.NewDeleteProductHandler(repo),
		updateStockHandler:          command.NewUpdateStockHandler(repo),
		reduceStockHandler:          command.NewReduceStockHandler(repo),
		increaseStockHandler:        command.NewIncreaseStockHandler(repo),
		activateProductHandler:      command.NewActivateProductHandler(repo),
		deactivateProductHandler:    command.NewDeactivateProductHandler(repo),
		markAsFeaturedHandler:       command.NewMarkAsFeaturedHandler(repo),
		unmarkAsFeaturedHandler:     command.NewUnmarkAsFeaturedHandler(repo),
		incrementViewCountHandler:   command.NewIncrementViewCountHandler(repo),

		// Query handlers
		getProductByIDHandler:       query.NewGetProductByIDHandler(repo),
		getProductBySKUHandler:      query.NewGetProductBySKUHandler(repo),
		listProductsHandler:         query.NewListProductsHandler(repo),
		listProductsByCategoryHandler: query.NewListProductsByCategoryHandler(repo),
		searchProductsHandler:       query.NewSearchProductsHandler(repo),
	}
}
