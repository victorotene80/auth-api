package infrastructure

import "errors"

var (
	ErrConflict        = errors.New("session version conflict, optimistic lock failed")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)
