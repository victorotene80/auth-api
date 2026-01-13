package command

type RevokeSessionCommand struct {
	SessionID string
	UserID    string
	IPAddress string
	UserAgent string
	DeviceID  string
}
