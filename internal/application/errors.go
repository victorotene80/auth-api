package application

import "errors"

var (
	ErrHandlerExists   = errors.New("handler already registered for command")
	ErrHandlerNotFound = errors.New("handler not found for command")
	ErrNilCommand      = errors.New("command cannot be nil")
	ErrInvalidResult   = errors.New("invalid result type")
	ErrSessionInvalid  = errors.New("session invalid or expired")
)
