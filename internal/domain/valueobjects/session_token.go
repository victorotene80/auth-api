package valueobjects

import "errors"

type SessionTokenHash struct {
	value string
}

func NewSessionTokenHash(value string) (SessionTokenHash, error) {
	if value == "" {
		return SessionTokenHash{}, errors.New("token hash cannot be empty")
	}

	if len(value) != 64 {
		return SessionTokenHash{}, errors.New("invalid token hash length")
	}

	return SessionTokenHash{value: value}, nil
}

func (t SessionTokenHash) Value() string {
	return t.value
}

func (t SessionTokenHash) Equals(other SessionTokenHash) bool {
	return t.value == other.value
}
