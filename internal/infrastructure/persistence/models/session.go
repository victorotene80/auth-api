package models

import "time"

type Session struct {
	ID                string     `db:"id"`
	UserID            string     `db:"user_id"`
	Role              string     `db:"role"` // <- added
	TokenHash         string     `db:"token_hash"`
	PreviousTokenHash *string    `db:"previous_token_hash"`
	RotationID        *string    `db:"rotation_id"`
	IPAddress         string     `db:"ip_address"`
	DeviceID          string     `db:"device_id"`
	UserAgent         string     `db:"user_agent"`
	Status            string     `db:"status"` // ACTIVE, REVOKED, EXPIRED
	CreatedAt         time.Time  `db:"created_at"`
	LastSeenAt        time.Time  `db:"last_seen_at"`
	ExpiresAt         time.Time  `db:"expires_at"`
	RevokedAt         *time.Time `db:"revoked_at"`
	Version           int        `db:"version"` // optimistic locking
}
