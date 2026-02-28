package command

type CreateUserCommand struct {
	Email             string
	Password          string
	FirstName         string
	LastName          string
	MiddleName        *string
	Phone             *string
	Locale            *string
	Timezone          *string
	AcceptTerms       bool
	IPAddress         string
	UserAgent         string
	DeviceID          string
	RequestID         string
	DeviceFingerprint string
	DeviceName        string
}
