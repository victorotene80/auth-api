package command

type LockAccountCommand struct {
	UserID    string
	Reason    string
	UserAgent string
	IPAddress string
	DeviceID  string
}
