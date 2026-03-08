package types

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

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
	OccurredAt time.Time
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
	//occuredAt time.Time,
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
			OccurredAt: time.Now(),
		},
		nil,
	)
}