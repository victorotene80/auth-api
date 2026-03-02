package models

import "time"

type Session struct {
	ID               string     `db:"id"`
	UserID           string     `db:"user_id"`
	TokenHash        string     `db:"token_hash"`
	RefreshTokenHash *string    `db:"refresh_token_hash"`
	IPAddress        *string    `db:"ip_address"`
	UserAgent        *string    `db:"user_agent"`
	DeviceFingerprint *string   `db:"device_fingerprint"`
	DeviceName       *string    `db:"device_name"`
	CountryCode      *string    `db:"country_code"`
	City             *string    `db:"city"`
	IsMFAVerified    bool       `db:"is_mfa_verified"`
	ImpersonatedBy   *string    `db:"impersonated_by"`
	LastActiveAt     time.Time  `db:"last_active_at"`
	ExpiresAt        time.Time  `db:"expires_at"`
	RevokedAt        *time.Time `db:"revoked_at"`
	RevokeReason     *string    `db:"revoke_reason"`
	CreatedAt        time.Time  `db:"created_at"`
}