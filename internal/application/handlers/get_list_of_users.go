package handlers

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/application/query"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type ListUsersHandler struct {
	repo repository.UserRepository
}

func NewListUsersHandler(repo repository.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{
		repo: repo,
	}
}

func (h *ListUsersHandler) Handle(ctx context.Context, q query.ListUsersQuery) ([]*aggregates.UserAggregate, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}

	var roleVO *valueobjects.Role
	if q.Role != nil {
		role, err := valueobjects.NewRole(*q.Role)
		if err != nil {
			return nil, 0, err
		}
		roleVO = &role
	}

	users, total, err := h.repo.List(ctx, q.Page, q.PageSize, roleVO, q.IsActive)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
