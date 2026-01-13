package contracts

import (
	"context"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

// SessionResult represents the outcome of creating a session
type SessionResult struct {
	SessionID    string
	AccessToken  Token
	RefreshToken Token
	ExpiresAt    time.Time
}

type SessionService interface {
	Create(
		ctx context.Context,
		userID string,
		role valueobjects.Role,
		ipAddress string,
		userAgent string,
		deviceID string,
	) (SessionResult, error)

	/*Rotate(
		ctx context.Context,
		sessionID string,
		oldRefreshToken string,
	) (SessionResult, error)

	Revoke(
		ctx context.Context,
		sessionID string,
		reason string,
	) error

	ValidateAccess(
		ctx context.Context,
		accessToken string,
	) (sessionID, userID, role string, err error)

	ValidateRefresh(
		ctx context.Context,
		refreshToken string,
	) (sessionID, userID string, err error)

	ListUserSessions(
		ctx context.Context,
		userID string,
	) ([]SessionResult, error)*/
}
