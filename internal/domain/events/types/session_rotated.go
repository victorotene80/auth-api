package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type SessionRotatedPayload struct {
	OldSessionID string
	NewSessionID string
	UserID       string
}

const SessionRotatedEventName = "session.rotated"

func NewSessionRotatedEvent(
	oldSessionID string,
	newSessionID string,
	userID string,
) events.DomainEvent {

	return events.NewEvent(
		SessionRotatedEventName,
		userID,
		SessionRotatedPayload{
			OldSessionID: oldSessionID,
			NewSessionID: newSessionID,
			UserID:       userID,
		},
		nil,
	)
}
