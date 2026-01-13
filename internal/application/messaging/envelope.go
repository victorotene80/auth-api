package messaging

import (
	"encoding/json"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type Envelope struct {
	ID          string
	EventName   string
	AggregateID string
	Payload     []byte
	OccurredAt  time.Time
	Version     int
	Metadata    map[string]string
}

func ToEnvelope(e events.DomainEvent, metadata map[string]string) (Envelope, error) {
	b, err := json.Marshal(e.Payload())
	if err != nil {
		return Envelope{}, err
	}
	return Envelope{
		ID:          e.EventID(),
		EventName:   e.EventName(),
		AggregateID: e.AggregateID(),
		Payload:     b,
		OccurredAt:  e.OccurredAt(),
		Version:     e.Version(),
		Metadata:    metadata,
	}, nil
}
