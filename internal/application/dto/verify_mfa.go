package dto

type VerifyMFADTO struct {
	UserID string `json:"user_id"`
	Method string `json:"method"`
	Code   string `json:"code"`
}
