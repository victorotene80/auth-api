package services

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
)

type AccountLockService struct {
	policy policy.AccountLockPolicy
}

func NewAccountLockService(p policy.AccountLockPolicy) *AccountLockService {
	return &AccountLockService{policy: p}
}

func (s *AccountLockService) ShouldLock(failedAttempts int, lastFailedAt time.Time) bool {
	if failedAttempts >= s.policy.MaxFailedAttempts {
		if time.Since(lastFailedAt) <= s.policy.ResetWindow {
			return true
		}
	}
	return false
}

func (s *AccountLockService) IsLocked(lockedAt time.Time) bool {
	return time.Since(lockedAt) < s.policy.LockDuration
}
