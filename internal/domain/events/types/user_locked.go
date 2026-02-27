package types

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type UserLockedPayload struct {
	UserID string
	Email  string
	Time   time.Time
}

const UserLockedEventName = "user.locked"

func NewUserLockedEvent(
	userID string,
	email string,
	timestamp time.Time,
) events.DomainEvent {
	return events.NewEvent(
		UserLockedEventName,
		userID,
		UserLockedPayload{
			UserID: userID,
			Email:  email,
			Time:   timestamp,
		},
		nil,
	)
}

