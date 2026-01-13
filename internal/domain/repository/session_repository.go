package repository

import (
	"context"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type SessionRepository interface {
	Save(ctx context.Context, session *aggregates.SessionAggregate) error
	FindActiveByKeyHash(
		ctx context.Context,
		tokenHash valueobjects.SessionTokenHash,
		now time.Time,
	) (*aggregates.SessionAggregate, error)
	RevokeByID(ctx context.Context, sessionID string, now time.Time) (*aggregates.SessionAggregate, error)
	RevokeAllForUser(ctx context.Context, userID string, now time.Time) ([]*aggregates.SessionAggregate, error)
	FindByID(ctx context.Context, sessionID string) (*aggregates.SessionAggregate, error)
	FindByTokenHash(
		ctx context.Context,
		tokenHash valueobjects.SessionTokenHash,
		now time.Time,
	) (*aggregates.SessionAggregate, error)
	RotateSessionToken(
		ctx context.Context,
		sessionID string,
		newToken valueobjects.SessionTokenHash,
		rotationID string,
		now time.Time,
	) (*aggregates.SessionAggregate, error)
}
