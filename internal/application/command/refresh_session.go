package command

type RefreshSessionCommand struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
	DeviceID     string
}
