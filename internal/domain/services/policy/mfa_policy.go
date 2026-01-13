package policy

import "time"

type MFAMethod string

const (
	MFAMethodTOTP  MFAMethod = "totp"
	MFAMethodSMS   MFAMethod = "sms"
	MFAMethodEmail MFAMethod = "email"
)

type MFAPolicy struct {
	Required         bool
	AllowedMethods   []MFAMethod
	GracePeriod      time.Duration
	RememberDevice   bool
	RememberDuration time.Duration
	RequiredActions  []string
}

func DefaultMFAPolicy() MFAPolicy {
	return MFAPolicy{
		Required:         false,
		AllowedMethods:   []MFAMethod{MFAMethodTOTP, MFAMethodEmail},
		GracePeriod:      7 * 24 * time.Hour,
		RememberDevice:   true,
		RememberDuration: 30 * 24 * time.Hour,
		RequiredActions:  []string{"change_password", "delete_account"},
	}
}
