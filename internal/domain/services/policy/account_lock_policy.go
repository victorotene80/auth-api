package policy

import "time"

type AccountLockPolicy struct {
	MaxFailedAttempts int
	LockDuration      time.Duration
	ResetWindow       time.Duration 
}

func DefaultAccountLockPolicy() AccountLockPolicy {
	return AccountLockPolicy{
		MaxFailedAttempts: 5,
		LockDuration:      15 * time.Minute,
		ResetWindow:       1 * time.Hour,
	}
}

func (p AccountLockPolicy) ShouldLockAccount(failedAttempts int, lastFailedAt time.Time) bool {
	if failedAttempts >= p.MaxFailedAttempts {
		if time.Since(lastFailedAt) <= p.ResetWindow {
			return true
		}
	}
	return false
}

func (p AccountLockPolicy) IsAccountStillLocked(lockedAt time.Time) bool {
	return time.Since(lockedAt) < p.LockDuration
}
