// internal/infrastructure/persistence/cache/session_cached.go
package cache

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
)

type CachedSession struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	TokenHash         string     `json:"token_hash"`
	IPAddress         string     `json:"ip_address,omitempty"`
	DeviceFingerprint string     `json:"device_fingerprint,omitempty"`
	DeviceName        string     `json:"device_name,omitempty"`
	UserAgent         string     `json:"user_agent,omitempty"`
	CountryCode       string     `json:"country_code,omitempty"`
	City              string     `json:"city,omitempty"`
	IsMFAVerified     bool       `json:"is_mfa_verified"`
	ImpersonatedBy    *string    `json:"impersonated_by,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	LastActiveAt      time.Time  `json:"last_active_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty"`
	RevokeReason      *string    `json:"revoke_reason,omitempty"`
	Version           int        `json:"version"`
}


func MapAggregateToCached(s *aggregates.SessionAggregate) CachedSession {
	return CachedSession{
		ID:                s.ID(),
		UserID:            s.UserID(),
		TokenHash:         s.TokenHash().Value(),
		IPAddress:         s.IPAddress(),
		DeviceFingerprint: s.DeviceFingerprint(),
		DeviceName:        s.DeviceName(),
		UserAgent:         s.UserAgent(),
		CountryCode:       s.CountryCode(),
		City:              s.City(),
		IsMFAVerified:     s.IsMFAVerified(),
		ImpersonatedBy:    s.ImpersonatedBy(),
		CreatedAt:         s.CreatedAt(),
		LastActiveAt:      s.LastActiveAt(),
		ExpiresAt:         s.ExpiresAt(),
		RevokedAt:         s.RevokedAt(),
		RevokeReason:      s.RevokeReason(),
		Version:           s.Version(),
	}
}

func MapCachedToAggregate(c CachedSession) (*aggregates.SessionAggregate, error) {
	return aggregates.RehydrateSession(
		c.ID,
		c.UserID,
		c.TokenHash,
		nil, // refreshTokenHash is not stored in cache
		c.IPAddress,
		c.UserAgent,
		c.DeviceFingerprint,
		c.DeviceName,
		c.CountryCode,
		c.City,
		c.IsMFAVerified,
		c.ImpersonatedBy,
		c.CreatedAt,
		c.LastActiveAt,
		c.ExpiresAt,
		c.RevokedAt,
		c.RevokeReason,
		c.Version,
	)
}