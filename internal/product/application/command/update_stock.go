package command

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// UpdateStockCommand represents the command to update product stock
type UpdateStockCommand struct {
	ProductID uint `json:"product_id"`
	Stock     int  `json:"stock"`
}

// UpdateStockHandler handles the update stock command
type UpdateStockHandler struct {
	repo domain.ProductRepository
}

// NewUpdateStockHandler creates a new update stock handler
func NewUpdateStockHandler(repo domain.ProductRepository) *UpdateStockHandler {
	return &UpdateStockHandler{
		repo: repo,
	}
}

// Handle executes the update stock command
func (h *UpdateStockHandler) Handle(ctx context.Context, cmd UpdateStockCommand) error {
	return h.repo.UpdateStock(ctx, cmd.ProductID, cmd.Stock)
}

// ReduceStockCommand represents the command to reduce product stock
type ReduceStockCommand struct {
	ProductID uint `json:"product_id"`
	Amount    int  `json:"amount"`
}

// ReduceStockHandler handles the reduce stock command
type ReduceStockHandler struct {
	repo domain.ProductRepository
}

// NewReduceStockHandler creates a new reduce stock handler
func NewReduceStockHandler(repo domain.ProductRepository) *ReduceStockHandler {
	return &ReduceStockHandler{
		repo: repo,
	}
}

// Handle executes the reduce stock command
func (h *ReduceStockHandler) Handle(ctx context.Context, cmd ReduceStockCommand) error {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	if err := product.ReduceStock(cmd.Amount); err != nil {
		return err
	}

	return h.repo.Update(ctx, product)
}

// IncreaseStockCommand represents the command to increase product stock
type IncreaseStockCommand struct {
	ProductID uint `json:"product_id"`
	Amount    int  `json:"amount"`
}

// IncreaseStockHandler handles the increase stock command
type IncreaseStockHandler struct {
	repo domain.ProductRepository
}

// NewIncreaseStockHandler creates a new increase stock handler
func NewIncreaseStockHandler(repo domain.ProductRepository) *IncreaseStockHandler {
	return &IncreaseStockHandler{
		repo: repo,
	}
}

// Handle executes the increase stock command
func (h *IncreaseStockHandler) Handle(ctx context.Context, cmd IncreaseStockCommand) error {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	if err := product.IncreaseStock(cmd.Amount); err != nil {
		return err
	}

	return h.repo.Update(ctx, product)
}
