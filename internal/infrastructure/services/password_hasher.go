package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	pepper []byte
	cost   int
}

func NewBcryptPasswordHasher(pepper string, cost int) (*BcryptPasswordHasher, error) {
	if len(pepper) == 0 {
		return nil, errors.New("pepper must not be empty")
	}

	return &BcryptPasswordHasher{
		pepper: []byte(pepper),
		cost:   cost,
	}, nil
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	mac := hmac.New(sha256.New, h.pepper)
	mac.Write([]byte(password))
	derived := mac.Sum(nil)

	hashed, err := bcrypt.GenerateFromPassword(derived, h.cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (h *BcryptPasswordHasher) Verify(plain, hash string) bool {
	mac := hmac.New(sha256.New, h.pepper)
	mac.Write([]byte(plain))
	derived := mac.Sum(nil)

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		derived,
	)

	return err == nil
}
