package outbox

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/events"
)

var _ contracts.EventPublisher = (*Publisher)(nil) // 👈 compile-time guarantee

type Publisher struct {
	repo       OutboxRepository
	serializer JSONSerializer
}

func NewPublisher(repo OutboxRepository) *Publisher {
	return &Publisher{
		repo:       repo,
		serializer: JSONSerializer{},
	}
}

func (p *Publisher) Publish(
	ctx context.Context,
	messages []events.DomainEvent,
	metadata map[string]string,
) error {

	for _, m := range messages {

		msg, err := messaging.ToEnvelope(m, metadata)
		if err != nil {
			return err
		}

		if err := p.repo.Add(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}
