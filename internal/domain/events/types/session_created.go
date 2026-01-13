package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type SessionCreatedPayload struct {
	SessionID string
	UserID    string
	IPAddress string
	UserAgent string
	ExpiresAt string
}

const SessionCreatedEventName = "session.created"

func NewSessionCreatedEvent(
	sessionID string,
	userID string,
	ipAddress string,
	userAgent string,
	expiresAt string,
) events.DomainEvent {

	return events.NewEvent(
		SessionCreatedEventName,
		userID,
		SessionCreatedPayload{
			SessionID: sessionID,
			UserID:    userID,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			ExpiresAt: expiresAt,
		},
		nil,
	)
}
