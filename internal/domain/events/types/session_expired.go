package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type SessionExpiredPayload struct {
	SessionID string
	UserID    string
}

const SessionExpiredEventName = "session.expired"

func NewSessionExpiredEvent(
	sessionID string,
	userID string,
) events.DomainEvent {

	return events.NewEvent(
		SessionExpiredEventName,
		userID,
		SessionExpiredPayload{
			SessionID: sessionID,
			UserID:    userID,
		},
		nil,
	)
}
