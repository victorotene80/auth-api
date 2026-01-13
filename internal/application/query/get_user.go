package query

type GetUserQuery struct {
	ID        *string
	Email     *string
	UserAgent string
	IPAddress string
	DeviceID  string
}
