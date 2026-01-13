package contracts

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/domain/events"
)

type EventPublisher interface {
	Publish(ctx context.Context, events []events.DomainEvent, metadata map[string]string) error
}
