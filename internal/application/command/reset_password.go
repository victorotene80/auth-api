package command

type ResetPasswordCommand struct {
	Email       string
	ResetToken  string
	NewPassword string
	IPAddress   string
	UserAgent   string
	DeviceID    string
}
