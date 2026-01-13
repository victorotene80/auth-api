package events

import "time"

type baseDomainEvent struct {
	id          string
	name        string
	occurredAt  time.Time
	aggregateID string
	payload     any
	meta        map[string]string
	version     int
}

func (e baseDomainEvent) EventID() string         { return e.id }
func (e baseDomainEvent) EventName() string       { return e.name }
func (e baseDomainEvent) AggregateID() string     { return e.aggregateID }
func (e baseDomainEvent) OccurredAt() time.Time   { return e.occurredAt }
func (e baseDomainEvent) Payload() any            { return e.payload }
func (e baseDomainEvent) Metadata() map[string]string { return e.meta }
func (e baseDomainEvent) Version() int            { return e.version }
