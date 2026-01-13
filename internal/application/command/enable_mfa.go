package command

type EnableMFACommand struct {
	UserID    string
	Method    string // "totp", "sms", "email"
	UserAgent string
	IPAddress string
	DeviceID  string
}
