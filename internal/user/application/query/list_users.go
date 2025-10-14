package query

import (
	"context"

	"github.com/ddd-micro/internal/user/domain"
)

// ListUsersQuery represents a query to list users with pagination
type ListUsersQuery struct {
	Offset int
	Limit  int
}

// ListUsersResult represents the result of listing users
type ListUsersResult struct {
	Users  []*domain.User
	Total  int
	Offset int
	Limit  int
}

// ListUsersHandler handles the ListUsersQuery
type ListUsersHandler struct {
	repo domain.UserRepository
}

// NewListUsersHandler creates a new ListUsersHandler
func NewListUsersHandler(repo domain.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{
		repo: repo,
	}
}

// Handle executes the ListUsersQuery
func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
	users, err := h.repo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	return &ListUsersResult{
		Users:  users,
		Total:  len(users),
		Offset: query.Offset,
		Limit:  query.Limit,
	}, nil
}

