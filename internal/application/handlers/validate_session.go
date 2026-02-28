package handlers

/*import (
	"context"
	"time"

	"github.com/victorotene80/authentication_api/internal/application"
	"github.com/victorotene80/authentication_api/internal/application/command"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type ValidateSessionHandler struct {
	uow       contracts.UnitOfWork
	repo      repository.SessionRepository
	publisher contracts.EventPublisher
	hasher    *utils.SessionKeyHasher
	clock     func() time.Time
}

func NewValidateSessionHandler(
	uow contracts.UnitOfWork,
	repo repository.SessionRepository,
	publisher contracts.EventPublisher,
	hasher *utils.SessionKeyHasher,
	clock func() time.Time,
) *ValidateSessionHandler {
	return &ValidateSessionHandler{
		uow:       uow,
		repo:      repo,
		publisher: publisher,
		hasher:    hasher,
		clock:     clock,
	}
}

func (h *ValidateSessionHandler) Handle(
	ctx context.Context,
	cmd command.ValidateSessionCommand,
) (*aggregates.SessionAggregate, error) {

	now := h.clock()

	hash := h.hasher.Hash(cmd.Token)
	tokenHash, err := valueobjects.NewSessionTokenHash(hash)
	if err != nil {
		return nil, err
	}

	var session *aggregates.SessionAggregate

	err = h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		s, err := h.repo.FindByTokenHash(txCtx, tokenHash, now)
		if err != nil {
			return err
		}

		s.EvaluateExpiry(now)

		if !s.IsValid(now) {
			return application.ErrSessionInvalid
		}

		s.Touch(now)

		if err := h.repo.Save(txCtx, s); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "session",
			Action:    "validated",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(txCtx, s.PullEvents(), meta.ToMetadata()); err != nil {
			return err
		}

		session = s
		return nil
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}*/
