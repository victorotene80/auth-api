package dto

type LockAccountDTO struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}