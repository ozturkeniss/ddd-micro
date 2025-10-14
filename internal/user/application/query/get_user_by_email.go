package query

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// GetUserByEmailQuery represents a query to get a user by email
type GetUserByEmailQuery struct {
	Email string
}

// GetUserByEmailHandler handles the GetUserByEmailQuery
type GetUserByEmailHandler struct {
	repo domain.UserRepository
}

// NewGetUserByEmailHandler creates a new GetUserByEmailHandler
func NewGetUserByEmailHandler(repo domain.UserRepository) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{
		repo: repo,
	}
}

// Handle executes the GetUserByEmailQuery
func (h *GetUserByEmailHandler) Handle(ctx context.Context, query GetUserByEmailQuery) (*domain.User, error) {
	return h.repo.GetByEmail(ctx, query.Email)
}

