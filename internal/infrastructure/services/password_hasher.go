package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost   int
	pepper string
}

func NewBcryptPasswordHasher(pepper string, cost int) (*BcryptPasswordHasher, error) {
	if pepper == "" {
		return nil, errors.New("password pepper must not be empty")
	}

	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptPasswordHasher{
		cost:   cost,
		pepper: pepper,
	}, nil
}

func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	// password + pepper (pepper is NOT stored)
	peppered := password + h.pepper

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(peppered),
		h.cost,
	)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (h *BcryptPasswordHasher) Verify(plain string, hash string) bool {
	peppered := plain + h.pepper

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(peppered),
	)

	return err == nil
}
