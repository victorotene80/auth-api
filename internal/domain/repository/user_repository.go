package repository

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type UserRepository interface {
	Create(ctx context.Context, user *aggregates.UserAggregate) error
	FindByID(ctx context.Context, id string) (*aggregates.UserAggregate, error)
	FindByEmail(ctx context.Context, email valueobjects.Email) (*aggregates.UserAggregate, error)
	ExistsByEmail(ctx context.Context, email valueobjects.Email) (bool, error)
	Update(ctx context.Context, user *aggregates.UserAggregate) error
	SoftDelete(ctx context.Context, id string) error
	List(
		ctx context.Context,
		limit int,
		cursor *string,
		status *valueobjects.UserStatus,
	) (users []*aggregates.UserAggregate, nextCursor *string, err error)
}
