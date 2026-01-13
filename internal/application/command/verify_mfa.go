package command

type VerifyMFACommand struct {
	UserID string
	Code   string
	Method string
}