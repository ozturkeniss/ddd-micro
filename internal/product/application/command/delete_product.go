package command

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// DeleteProductCommand represents the command to delete a product
type DeleteProductCommand struct {
	ProductID uint `json:"product_id"`
}

// DeleteProductHandler handles the delete product command
type DeleteProductHandler struct {
	repo domain.ProductRepository
}

// NewDeleteProductHandler creates a new delete product handler
func NewDeleteProductHandler(repo domain.ProductRepository) *DeleteProductHandler {
	return &DeleteProductHandler{
		repo: repo,
	}
}

// Handle executes the delete product command
func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	return h.repo.Delete(ctx, cmd.ProductID)
}
