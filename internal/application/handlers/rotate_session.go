package handlers

/*import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type RotateSessionHandler struct {
	uow       appContracts.UnitOfWork
	repo      repository.SessionRepository
	publisher appContracts.EventPublisher
	hasher    *utils.SessionKeyHasher
	clock     func() time.Time
	policy    policy.SessionPolicy
	tokenGen  contracts.TokenGenerator
}

func NewRotateSessionHandler(
	uow appContracts.UnitOfWork,
	repo repository.SessionRepository,
	publisher appContracts.EventPublisher,
	hasher *utils.SessionKeyHasher,
	policy policy.SessionPolicy,
	tokenGen contracts.TokenGenerator,
	clock func() time.Time,
) *RotateSessionHandler {
	return &RotateSessionHandler{
		uow:       uow,
		repo:      repo,
		publisher: publisher,
		hasher:    hasher,
		clock:     clock,
		policy:    policy,
		tokenGen:  tokenGen,
	}
}

func (h *RotateSessionHandler) Handle(
	ctx context.Context,
	cmd command.RotateSessionCommand,
) (*dto.RefreshResultDTO, error) {

	now := h.clock()

	oldHash := h.hasher.Hash(cmd.RawToken)
	oldTokenHash, err := valueobjects.NewSessionTokenHash(oldHash)
	if err != nil {
		return nil, err
	}

	newRawRefreshToken, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	newTokenHash, err := valueobjects.NewSessionTokenHash(h.hasher.Hash(newRawRefreshToken))
	if err != nil {
		return nil, err
	}

	var result *dto.RefreshResultDTO

	err = h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		session, err := h.repo.FindByTokenHash(txCtx, oldTokenHash, now)
		if err != nil {
			return errors.New("invalid refresh token")
		}

		if session.Status() != aggregates.SessionActive {
			return errors.New("session is not active")
		}

		if session.PreviousTokenHash() != nil && session.PreviousTokenHash().Equals(oldTokenHash) {
			session.Revoke(now, "Refresh token replay detected")
			if err := h.repo.Save(txCtx, session); err != nil {
				return err
			}
			return errors.New("security violation: refresh token reused")
		}

		if !session.IsValid(now) {
			session.EvaluateExpiry(now)
			if err := h.repo.Save(txCtx, session); err != nil {
				return err
			}
			return errors.New("session expired")
		}

		rotationID := uuid.NewString()
		
		if err := session.RotateKey(newTokenHash, rotationID, now); err != nil {
			return err
		}

		if err := h.repo.Save(txCtx, session); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "session",
			Action:    "rotated",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}
		if err := h.publisher.Publish(txCtx, session.PullEvents(), meta.ToMetadata()); err != nil {
			return err
		}

		accessToken, err := h.tokenGen.GenerateAccess(session.UserID(),session.ID(), session.Role().String(), h.policy.MaxDuration)
		if err != nil {
			return err
		}

		result = &dto.RefreshResultDTO{
			AccessToken:  accessToken,
			RefreshToken: newRawRefreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
*/