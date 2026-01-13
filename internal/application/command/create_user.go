package command

type CreateUserCommand struct {
	Email      string
	Password   string
	FirstName  string
	LastName   string
	MiddleName string
	IPAddress  string
	UserAgent  string
	DeviceID   string
}
