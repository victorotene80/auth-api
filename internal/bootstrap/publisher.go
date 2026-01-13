package bootstrap

import (
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/infrastructure/messaging/outbox"
)

func initializeEventPublisher(p *Persistence) appContracts.EventPublisher {
	if p.OutboxRepo == nil {
		panic("Outbox repository is not initialized")
	}

	pub := outbox.NewPublisher(p.OutboxRepo)
	return pub
}
