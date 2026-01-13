package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type UserCreatedPayload struct {
	UserID    string
	Email     string
	Role      string
	FirstName string
	LastName  string
	//Status    string
}

const UserCreatedEventName = "user.created"

func NewUserCreatedEvent(
	userID string,
	email string,
	role string,
	firstName string,
	lastName string,
	//status string,
	//registrationType string,
) events.DomainEvent {

	return events.NewEvent(
		UserCreatedEventName,
		userID,
		UserCreatedPayload{
			UserID:    userID,
			Email:     email,
			Role:      role,
			FirstName: firstName,
			LastName:  lastName,
			//Status:    status,
		},
		nil,
	)
}
