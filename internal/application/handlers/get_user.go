package handlers

import (
	"context"
	"errors"

	"github.com/victorotene80/authentication_api/internal/application/query"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type GetUserHandler struct {
	repo repository.UserRepository
}

func NewGetUserHandler(repo repository.UserRepository) *GetUserHandler {
	return &GetUserHandler{
		repo: repo,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, q query.GetUserQuery) (*aggregates.UserAggregate, error) {
	if q.ID == nil && q.Email == nil {
		return nil, errors.New("either ID or Email must be provided")
	}

	if q.ID != nil {
		user, err := h.repo.FindByID(ctx, *q.ID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		return user, nil
	}

	emailVO, err := valueobjects.NewEmail(*q.Email)
	if err != nil {
		return nil, err
	}

	user, err := h.repo.FindByEmail(ctx, emailVO)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
