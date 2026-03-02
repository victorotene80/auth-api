package bootstrap

import (
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/infrastructure/messaging/outbox"
	"go.uber.org/zap"
)

func initializeEventPublisher(
	p *Persistence,
	logger *zap.Logger,
) appContracts.EventPublisher {

	if p.OutboxRepo == nil {
		logger.Fatal("Outbox repository is nil — cannot initialize event publisher")
	}

	logger.Info("Initializing event publisher")
	return outbox.NewPublisher(p.OutboxRepo)
}