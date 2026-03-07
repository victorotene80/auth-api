package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type UserCreatedPayload struct {
	UserID     string
	Email      string
	Phone      string
	//Role       string
	FirstName  string
	LastName   string
	Status     string
	MFAEnabled bool
	Version    int
}

const UserCreatedEventName = "user.created"

func NewUserCreatedEvent(
	userID string,
	email string,
	phone string,
	//role string,
	firstName string,
	lastName string,
	status string,
	mfaEnabled bool,
	version int,
) events.DomainEvent {
	return events.NewEvent(
		UserCreatedEventName,
		userID,
		UserCreatedPayload{
			UserID:     userID,
			Email:      email,
			Phone:      phone,
			//Role:       role,
			FirstName:  firstName,
			LastName:   lastName,
			Status:     status,
			MFAEnabled: mfaEnabled,
			Version:    version,
		},
		nil,
	)
}