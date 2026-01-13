package bootstrap

import (
	"database/sql"
	"fmt"

	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence"
	"github.com/victorotene80/authentication_api/internal/shared/config"
)

type Persistence struct {
	DB          *sql.DB
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
	OutboxRepo  contracts.OutboxRepository
	UoW         contracts.UnitOfWork
}

func initializePersistence(cfg *config.Config) (*Persistence, error) {
	db, err := persistence.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	return &Persistence{
		DB:          db,
		UserRepo:    persistence.NewPostgresUserRepository(db),
		SessionRepo: persistence.NewPostgresSessionRepository(db),
		OutboxRepo:  persistence.NewPostgresOutboxRepository(db),
		UoW:         persistence.NewSqlUnitOfWork(db),
	}, nil
}
//pa55w0rd5