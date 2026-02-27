package types

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type LoginFailedPayload struct {
	UserID string
	Email  string
	Time   time.Time
}

type UserLoginFailedEvent struct {
	UserID    string
	UserEmail string
	Timestamp time.Time
}

const UserLoginFailedEventName = "user.loginFailed"

func NewUserLoginFailedEvent(
	userID string,
	userEmail string,
	timestamp time.Time,
) events.DomainEvent {
	return events.NewEvent(
		UserLoginFailedEventName,
		userID,
		LoginFailedPayload{
			UserID:    userID,
			Email:     userEmail,
			Time:      timestamp,
		},
		nil,
	)
}

