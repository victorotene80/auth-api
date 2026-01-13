package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type SessionAccessedPayload struct {
	SessionID string
	UserID    string
}

const SessionAccessedEventName = "session.accessed"

func NewSessionAccessedEvent(
	sessionID string,
	userID string,
) events.DomainEvent {

	return events.NewEvent(
		SessionAccessedEventName,
		userID,
		SessionAccessedPayload{
			SessionID: sessionID,
			UserID:    userID,
		},
		nil,
	)
}
