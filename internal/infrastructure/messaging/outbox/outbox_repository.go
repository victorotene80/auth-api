package outbox

import (
	"context"

	"github.com/victorotene80/authentication_api/internal/application/messaging"
)

type OutboxRepository interface {
	Add(ctx context.Context, envelope messaging.Envelope) error
}
