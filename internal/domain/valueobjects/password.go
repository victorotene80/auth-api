package valueobjects

import (
	"github.com/victorotene80/authentication_api/internal/domain"
)

type Password struct {
	value string
}

func NewHashedPassword(hash string) (Password, error) {
	if hash == "" {
		return Password{}, domain.ErrInvalidPasswordHash
	}
	return Password{value: hash}, nil
}

func EmptyPassword() Password {
	return Password{value: ""}
}

func (p Password) Value() string {
	return p.value
}

func (p Password) IsEmpty() bool {
	return p.value == ""
}
