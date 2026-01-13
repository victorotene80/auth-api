package dto

type EnableMFADTO struct {
	UserID string `json:"user_id"`
	Method string `json:"method"` // totp, sms, email
}
