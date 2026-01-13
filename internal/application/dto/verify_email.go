package dto

type VerifyEmailDTO struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}
