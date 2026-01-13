package cache

import (
	"time"
)

type CachedSession struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	Role              string     `json:"role"`
	TokenHash         string     `json:"token_hash"`
	PreviousTokenHash *string    `json:"previous_token_hash,omitempty"`
	RotationID        *string    `json:"rotation_id,omitempty"`
	IPAddress         string     `json:"ip_address"`
	DeviceID          string     `json:"device_id"`
	UserAgent         string     `json:"user_agent"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	LastSeenAt        time.Time  `json:"last_seen_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty"`
	Version           int        `json:"version"`
}
