package bootstrap

import (
	"database/sql"
	"fmt"

	outboxContracts "github.com/victorotene80/authentication_api/internal/infrastructure/messaging/outbox"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence"
	"github.com/victorotene80/authentication_api/internal/shared/config"
	"go.uber.org/zap"
)

type Persistence struct {
	DB          *sql.DB
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
	OutboxRepo  outboxContracts.OutboxRepository
	RoleRepo    repository.RoleRepository
	UoW         contracts.UnitOfWork
}

func initializePersistence(cfg *config.Config, logger *zap.Logger) (*Persistence, error) {
	db, err := persistence.NewDatabase(cfg.Database)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
		return nil, fmt.Errorf("db init failed: %w", err)
	}

	logger.Info("database connected",
		zap.String("addr", fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port)),
		zap.String("db", cfg.Database.Name),
	)

	return &Persistence{
		DB:          db,
		UserRepo:    persistence.NewPostgresUserRepository(db),
		SessionRepo: persistence.NewPostgresSessionRepository(db),
		OutboxRepo:  persistence.NewPostgresOutboxRepository(db),
		RoleRepo:    persistence.NewPgRoleRepository(db),
		UoW:         persistence.NewSqlUnitOfWork(db),
	}, nil
}