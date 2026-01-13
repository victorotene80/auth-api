package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
)

// MFAService handles multi-factor authentication rules and challenges.
type MFAService struct {
	policy policy.MFAPolicy
}

// NewMFAService creates a new MFA service with a given policy.
func NewMFAService(p policy.MFAPolicy) *MFAService {
	return &MFAService{policy: p}
}

// IsMFARequiredForAction checks if MFA is required for a specific user action.
func (s *MFAService) IsMFARequiredForAction(action string) bool {
	if !s.policy.Required {
		return false
	}
	for _, a := range s.policy.RequiredActions {
		if a == action {
			return true
		}
	}
	return false
}

// IsWithinGracePeriod checks if the last MFA validation is still within the grace period.
func (s *MFAService) IsWithinGracePeriod(lastValidated, now time.Time) bool {
	return now.Sub(lastValidated) <= s.policy.GracePeriod
}

// IsDeviceRemembered checks if a device is remembered and still valid.
func (s *MFAService) IsDeviceRemembered(rememberedAt, now time.Time) bool {
	if !s.policy.RememberDevice {
		return false
	}
	return now.Sub(rememberedAt) <= s.policy.RememberDuration
}

// StartMFAChallenge generates a secure MFA challenge ID for a user action.
func (s *MFAService) StartMFAChallenge() (string, error) {
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		return "", errors.New("failed to generate MFA challenge")
	}
	return hex.EncodeToString(id), nil
}

// ValidateMFAMethod checks if the requested MFA method is allowed by policy.
func (s *MFAService) ValidateMFAMethod(method policy.MFAMethod) bool {
	for _, m := range s.policy.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}
