package command

import (
	"context"

	"github.com/ddd-micro/internal/product/domain"
)

// ActivateProductCommand represents the command to activate a product
type ActivateProductCommand struct {
	ProductID uint `json:"product_id"`
}

// ActivateProductHandler handles the activate product command
type ActivateProductHandler struct {
	repo domain.ProductRepository
}

// NewActivateProductHandler creates a new activate product handler
func NewActivateProductHandler(repo domain.ProductRepository) *ActivateProductHandler {
	return &ActivateProductHandler{
		repo: repo,
	}
}

// Handle executes the activate product command
func (h *ActivateProductHandler) Handle(ctx context.Context, cmd ActivateProductCommand) (*domain.Product, error) {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	product.Activate()

	if err := h.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeactivateProductCommand represents the command to deactivate a product
type DeactivateProductCommand struct {
	ProductID uint `json:"product_id"`
}

// DeactivateProductHandler handles the deactivate product command
type DeactivateProductHandler struct {
	repo domain.ProductRepository
}

// NewDeactivateProductHandler creates a new deactivate product handler
func NewDeactivateProductHandler(repo domain.ProductRepository) *DeactivateProductHandler {
	return &DeactivateProductHandler{
		repo: repo,
	}
}

// Handle executes the deactivate product command
func (h *DeactivateProductHandler) Handle(ctx context.Context, cmd DeactivateProductCommand) (*domain.Product, error) {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	product.Deactivate()

	if err := h.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// MarkAsFeaturedCommand represents the command to mark a product as featured
type MarkAsFeaturedCommand struct {
	ProductID uint `json:"product_id"`
}

// MarkAsFeaturedHandler handles the mark as featured command
type MarkAsFeaturedHandler struct {
	repo domain.ProductRepository
}

// NewMarkAsFeaturedHandler creates a new mark as featured handler
func NewMarkAsFeaturedHandler(repo domain.ProductRepository) *MarkAsFeaturedHandler {
	return &MarkAsFeaturedHandler{
		repo: repo,
	}
}

// Handle executes the mark as featured command
func (h *MarkAsFeaturedHandler) Handle(ctx context.Context, cmd MarkAsFeaturedCommand) (*domain.Product, error) {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	product.MarkAsFeatured()

	if err := h.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// UnmarkAsFeaturedCommand represents the command to unmark a product as featured
type UnmarkAsFeaturedCommand struct {
	ProductID uint `json:"product_id"`
}

// UnmarkAsFeaturedHandler handles the unmark as featured command
type UnmarkAsFeaturedHandler struct {
	repo domain.ProductRepository
}

// NewUnmarkAsFeaturedHandler creates a new unmark as featured handler
func NewUnmarkAsFeaturedHandler(repo domain.ProductRepository) *UnmarkAsFeaturedHandler {
	return &UnmarkAsFeaturedHandler{
		repo: repo,
	}
}

// Handle executes the unmark as featured command
func (h *UnmarkAsFeaturedHandler) Handle(ctx context.Context, cmd UnmarkAsFeaturedCommand) (*domain.Product, error) {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, err
	}

	product.UnmarkAsFeatured()

	if err := h.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// IncrementViewCountCommand represents the command to increment product view count
type IncrementViewCountCommand struct {
	ProductID uint `json:"product_id"`
}

// IncrementViewCountHandler handles the increment view count command
type IncrementViewCountHandler struct {
	repo domain.ProductRepository
}

// NewIncrementViewCountHandler creates a new increment view count handler
func NewIncrementViewCountHandler(repo domain.ProductRepository) *IncrementViewCountHandler {
	return &IncrementViewCountHandler{
		repo: repo,
	}
}

// Handle executes the increment view count command
func (h *IncrementViewCountHandler) Handle(ctx context.Context, cmd IncrementViewCountCommand) error {
	product, err := h.repo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	product.IncrementViewCount()

	return h.repo.Update(ctx, product)
}
