package domain

import "errors"

var (
	ErrEmptyEmail               = errors.New("")
	ErrInvalidEmailFormat       = errors.New("")
	ErrInvalidRole              = errors.New("invalid role")
	ErrInvalidPasswordHash      = errors.New("invalid password hash")
	ErrPasswordTooShort         = errors.New("password is too short")
	ErrPasswordMissingUppercase = errors.New("password must contain uppercase letters")
	ErrPasswordMissingLowercase = errors.New("password must contain lowercase letters")
	ErrPasswordMissingNumber    = errors.New("password must contain a number")
	ErrPasswordMissingSpecial   = errors.New("password must contain a special character")
	ErrEmailAlreadyExists       = errors.New("email already exists")
)
