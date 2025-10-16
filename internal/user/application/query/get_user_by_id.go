package query

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// GetUserByIDQuery represents a query to get a user by ID
type GetUserByIDQuery struct {
	UserID uint
}

// GetUserByIDHandler handles the GetUserByIDQuery
type GetUserByIDHandler struct {
	repo domain.UserRepository
}

// NewGetUserByIDHandler creates a new GetUserByIDHandler
func NewGetUserByIDHandler(repo domain.UserRepository) *GetUserByIDHandler {
	return &GetUserByIDHandler{
		repo: repo,
	}
}

// Handle executes the GetUserByIDQuery
func (h *GetUserByIDHandler) Handle(ctx context.Context, query GetUserByIDQuery) (*domain.User, error) {
	return h.repo.GetByID(ctx, query.UserID)
}
