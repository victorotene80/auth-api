package response

import "time"

type LoginResponse struct {
	Tokens      AuthTokensResponse `json:"tokens,omitempty"`
	MFARequired bool               `json:"mfa_required"`
	LastLogin   *time.Time         `json:"last_login,omitempty"`
	ChallengeID *string            `json:"challenge_id,omitempty"`
	Status      string             `json:"status"`
}
