package contracts

import (
	"time"

)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}

type TokenGenerator interface {
	GenerateAccess(
		userID, sessionID string,
		duration time.Duration,
	) (Token, error)

	GenerateRefresh(
		userID string,
		sessionID string,
		duration time.Duration,
	) (Token, error)

	ValidateAccess(
		token string,
	) (userID string, sessionID string, err error)

	ValidateRefresh(
		token string,
	) (userID string, sessionID string, err error)
}
