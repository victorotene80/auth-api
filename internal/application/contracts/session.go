package contracts

import (
	"context"
	"time"

	domainToken "github.com/victorotene80/authentication_api/internal/domain/contracts"
)


type SessionResult struct {
	SessionID    string
	AccessToken  domainToken.Token
	RefreshToken domainToken.Token
	ExpiresAt    time.Time 
}


type SessionService interface {
	Create(
		ctx context.Context,
		userID string,
		ipAddress string,
		userAgent string,
		deviceID string,
		deviceFingerprint string,
		deviceName string,
	) (SessionResult, error)
}