package command

type ChangePasswordCommand struct {
	UserID      string
	OldPassword string
	NewPassword string
	UserAgent   string
	IPAddress   string
	DeviceID    string
}
