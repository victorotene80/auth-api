package request

type EnableMFARequest struct {
	Method string `json:"method" binding:"required,oneof=totp sms email"`
}
