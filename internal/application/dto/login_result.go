package dto

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/contracts"
)

type LoginResultDTO struct {
	Tokens      contracts.TokenPair
	MFARequired bool
	LastLogin   *time.Time
	ChallengeID *string
	Status string
}
