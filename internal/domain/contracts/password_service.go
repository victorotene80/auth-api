package contracts

type PasswordService interface {
	Validate(password string) error
}
