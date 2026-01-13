package command

type ValidateSessionCommand struct {
	Token string
	IPAddress    string
	UserAgent    string
	DeviceID     string
}