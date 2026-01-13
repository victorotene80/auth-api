package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type UserPasswordChangedPayload struct {
	UserID string
	Email  string
}

const UserPasswordChangedEventName = "user.password.changed"

func NewUserPasswordChangedEvent(
	userID string,
	email string,
) events.DomainEvent {

	return events.NewEvent(
		UserPasswordChangedEventName,
		userID,
		UserPasswordChangedPayload{
			UserID: userID,
			Email:  email,
		},
		nil,
	)
}
