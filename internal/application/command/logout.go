package command

type LogoutCommand struct {
	SessionID string
	UserID    string
	Reason    string // optional, e.g., "user logout"
	IPAddress string
	UserAgent string
	DeviceID  string
}
