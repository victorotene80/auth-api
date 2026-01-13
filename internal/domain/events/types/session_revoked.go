package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type SessionRevokedPayload struct {
	SessionID string
	UserID    string
	Reason    string
}

const SessionRevokedEventName = "session.revoked"

func NewSessionRevokedEvent(
	sessionID string,
	userID string,
	reason string,
) events.DomainEvent {

	return events.NewEvent(
		SessionRevokedEventName,
		userID,
		SessionRevokedPayload{
			SessionID: sessionID,
			UserID:    userID,
			Reason:    reason,
		},
		nil,
	)
}
