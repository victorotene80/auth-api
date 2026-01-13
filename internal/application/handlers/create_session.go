package handlers

/*import (
	"context"

	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type CreateSessionHandler struct {
	uow       contracts.UnitOfWork
	repo      repository.SessionRepository
	publisher contracts.EventPublisher
	hasher    *utils.SessionKeyHasher
	policy    policy.SessionPolicy
	clock     func() time.Time
}

func NewCreateSessionHandler(
	uow contracts.UnitOfWork,
	repo repository.SessionRepository,
	publisher contracts.EventPublisher,
	hasher *utils.SessionKeyHasher,
	policy policy.SessionPolicy,
	clock func() time.Time,
) *CreateSessionHandler {
	return &CreateSessionHandler{
		uow:       uow,
		repo:      repo,
		publisher: publisher,
		hasher:    hasher,
		policy:    policy,
		clock:     clock,
	}
}

func (h *CreateSessionHandler) Handle(
	ctx context.Context,
	cmd command.CreateSessionCommand,
) (rawToken string, err error) {

	now := h.clock()

	rawToken, err = utils.GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	hash := h.hasher.Hash(rawToken)
	tokenHash, err := valueobjects.NewSessionTokenHash(hash)
	if err != nil {
		return "", err
	}

	expiresAt := now.Add(h.policy.MaxDuration)

	err = h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {

		session, err := aggregates.NewSession(
			cmd.UserID,
			
			tokenHash,
			cmd.IPAddress,
			cmd.UserAgent,
			cmd.DeviceID,
			now,
			expiresAt,
		)
		if err != nil {
			return err
		}

		if err := h.repo.Save(txCtx, session); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "user",
			Action:    "created",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(
			ctx,
			session.PullEvents(),
			meta.ToMetadata(),
		); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return rawToken, nil
}
*/