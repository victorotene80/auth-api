package policy

import "time"

type SessionPolicy struct {
	MaxDuration           time.Duration
	IdleTimeout           time.Duration
	MaxConcurrentSessions int
	RequireMFAFor         []string
	AllowRefreshToken     bool
	RefreshTokenDuration  time.Duration
}

func DefaultSessionPolicy() SessionPolicy {
	return SessionPolicy{
		MaxDuration:           24 * time.Hour,
		IdleTimeout:           30 * time.Minute,
		MaxConcurrentSessions: 5,
		RequireMFAFor:         []string{"change_password", "delete_account"},
		AllowRefreshToken:     true,
		RefreshTokenDuration:  7 * 24 * time.Hour,
	}
}

func (p SessionPolicy) IsSessionExpired(createdAt, lastSeen, now time.Time) bool {
	if now.Sub(createdAt) > p.MaxDuration {
		return true
	}
	if now.Sub(lastSeen) > p.IdleTimeout {
		return true
	}
	return false
}

func (p SessionPolicy) ComputeExpiresAt(now time.Time) time.Time {
	return now.Add(p.MaxDuration)
}
