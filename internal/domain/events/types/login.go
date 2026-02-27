package types

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type UserLoginPayload struct {
	UserID    string
	Email    string
	Time      time.Time
}

const UserLoginEventName = "user.loggedIn"

func NewUserLogInEvent(
	userID string,
	email string,
	timestamp time.Time,
) events.DomainEvent {

	return events.NewEvent(
		UserLoginEventName,
		userID,
		UserLoginPayload{
			UserID:    userID,
			Email:    email,
			Time:      timestamp,
		},
		nil,
	)
}
