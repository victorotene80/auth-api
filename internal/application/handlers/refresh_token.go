package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/victorotene80/authentication_api/internal/application/command"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type RefreshSessionHandler struct {
	uow       contracts.UnitOfWork
	repo      repository.SessionRepository
	publisher contracts.EventPublisher
	hasher    *utils.SessionKeyHasher
	policy    policy.SessionPolicy
	clock     func() time.Time
}

func NewRefreshSessionHandler(
	uow contracts.UnitOfWork,
	repo repository.SessionRepository,
	publisher contracts.EventPublisher,
	hasher *utils.SessionKeyHasher,
	policy policy.SessionPolicy,
	clock func() time.Time,
) *RefreshSessionHandler {
	return &RefreshSessionHandler{
		uow:       uow,
		repo:      repo,
		publisher: publisher,
		hasher:    hasher,
		policy:    policy,
		clock:     clock,
	}
}

func (h *RefreshSessionHandler) Handle(
	ctx context.Context,
	cmd command.RefreshSessionCommand,
) (newRawToken string, err error) {

	now := h.clock()

	tokenHash := h.hasher.Hash(cmd.RefreshToken)

	tokenHashVO, err := valueobjects.NewSessionTokenHash(tokenHash)
	if err != nil {
		return "", err
	}

	err = h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		session, err := h.repo.FindByTokenHash(txCtx, tokenHashVO, now)
		if err != nil {
			return errors.New("invalid refresh token")
		}

		if session.Status() != aggregates.SessionActive {
			return errors.New("session is not active")
		}

		if h.policy.IsSessionExpired(session.CreatedAt(), session.LastSeenAt()) {
			session.EvaluateExpiry(now)

			if err := h.repo.Save(txCtx, session); err != nil {
				return err
			}

			return errors.New("session expired")
		}

		newRawToken, err = utils.GenerateRandomString(32)
		if err != nil {
			return err
		}

		newTokenHash, err := valueobjects.NewSessionTokenHash(h.hasher.Hash(newRawToken))
		if err != nil {
			return err
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
			Action:    "refreshed",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(txCtx, session.PullEvents(), meta.ToMetadata()); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newRawToken, nil
}
