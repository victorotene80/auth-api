package contracts

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(plain string, hash string) bool
}
