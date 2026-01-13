package events

import "time"

type DomainEvent interface {
	EventID() string
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
	Version() int
	Payload() any
	Metadata() map[string]string
}

