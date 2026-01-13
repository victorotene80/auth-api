package contracts // application layer

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
)

type SessionCache interface {
	Get(ctx context.Context, token string) (*aggregates.SessionAggregate, error)
	Set(ctx context.Context, token string, session *aggregates.SessionAggregate) error
	Delete(ctx context.Context, token string) error
	RefreshTTL(ctx context.Context, token string) error
}
