package aggregates

import (
	"github.com/victorotene80/authentication_api/internal/domain/events"
	"github.com/victorotene80/authentication_api/internal/shared/utils"

	"time"
)

type AggregateRoot struct {
	id                string
	version           int
	createdAt         time.Time
	updatedAt         time.Time
	uncommittedEvents []events.DomainEvent
}

func NewAggregateRoot(id string) *AggregateRoot {
	now := utils.NowUTC()
	return &AggregateRoot{
		id:                id,
		version:           0,
		createdAt:         now,
		updatedAt:         now,
		uncommittedEvents: make([]events.DomainEvent, 0),
	}
}

func (a *AggregateRoot) ID() string {
	return a.id
}

func (a *AggregateRoot) Version() int {
	return a.version
}

func (a *AggregateRoot) CreatedAt() time.Time {
	return a.createdAt
}

func (a *AggregateRoot) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *AggregateRoot) RaiseEvent(event events.DomainEvent) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
	a.updatedAt = utils.NowUTC()
	//a.incrementVersion()
}

func (a *AggregateRoot) PullEvents() []events.DomainEvent {
	eventsCopy := make([]events.DomainEvent, len(a.uncommittedEvents))
	copy(eventsCopy, a.uncommittedEvents)
	return eventsCopy
}

func (a *AggregateRoot) ClearEvents() {
	a.uncommittedEvents = make([]events.DomainEvent, 0)
}

//func (a *AggregateRoot) incrementVersion() {
//	a.version++
//}

func (a *AggregateRoot) CommitVersion() {
	a.version++
}
