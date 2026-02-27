package aggregates

import (
	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type AggregateRoot struct {
	id                string
	version           int 
	uncommittedEvents []events.DomainEvent
}

func NewAggregateRoot(id string, version int) *AggregateRoot {
	return &AggregateRoot{
		id:                id,
		version:           version,
		uncommittedEvents: make([]events.DomainEvent, 0),
	}
}

func (a *AggregateRoot) ID() string    { return a.id }
func (a *AggregateRoot) Version() int  { return a.version }
func (a *AggregateRoot) SetVersion(v int) {
	a.version = v
}

func (a *AggregateRoot) RaiseEvent(event events.DomainEvent) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
}

func (a *AggregateRoot) PullEvents() []events.DomainEvent {
	eventsCopy := make([]events.DomainEvent, len(a.uncommittedEvents))
	copy(eventsCopy, a.uncommittedEvents)
	return eventsCopy
}

func (a *AggregateRoot) ClearEvents() {
	a.uncommittedEvents = a.uncommittedEvents[:0]
}

func (a *AggregateRoot) CommitVersion(newVersion int) {
	a.version = newVersion
}