package valueobjects

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type PhoneNumber struct {
	value string
}

var phoneDigitsOnlyRegex = regexp.MustCompile(`^\+?[1-9][0-9]{7,14}$`)

func NewPhoneNumber(input string) (PhoneNumber, error) {
	normalized := normalizePhoneNumber(input)

	if normalized == "" {
		return PhoneNumber{}, errors.New("phone number is required")
	}

	if !phoneDigitsOnlyRegex.MatchString(normalized) {
		return PhoneNumber{}, fmt.Errorf("invalid phone number format: %q", input)
	}

	return PhoneNumber{value: normalized}, nil
}

func MustNewPhoneNumber(input string) PhoneNumber {
	p, err := NewPhoneNumber(input)
	if err != nil {
		panic(err)
	}
	return p
}

func (p PhoneNumber) String() string {
	return p.value
}

func (p PhoneNumber) Value() string {
	return p.value
}

func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

func (p PhoneNumber) IsZero() bool {
	return p.value == ""
}

func normalizePhoneNumber(input string) string {
	s := strings.TrimSpace(input)

	replacer := strings.NewReplacer(
		" ", "",
		"-", "",
		"(", "",
		")", "",
	)
	s = replacer.Replace(s)

	if strings.HasPrefix(s, "00") {
		s = "+" + strings.TrimPrefix(s, "00")
	}

	return s
}
