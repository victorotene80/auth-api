package handlers

import (
	"context"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	"github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
)

type RevokeSessionHandler struct {
	uow       contracts.UnitOfWork
	repo      repository.SessionRepository
	publisher contracts.EventPublisher
	clock     func() time.Time
}

func NewRevokeSessionHandler(
	uow contracts.UnitOfWork,
	repo repository.SessionRepository,
	publisher contracts.EventPublisher,
	clock func() time.Time,
) *RevokeSessionHandler {
	return &RevokeSessionHandler{
		uow:       uow,
		repo:      repo,
		publisher: publisher,
		clock:     clock,
	}
}

func (h *RevokeSessionHandler) Handle(
	ctx context.Context,
	cmd command.RevokeSessionCommand,
) error {
	now := h.clock()

	return h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		session, err := h.repo.FindByID(txCtx, cmd.SessionID)
		if err != nil {
			return err
		}

		session.EvaluateExpiry(now)

		if err := session.Revoke(now, "User initiated: "+cmd.UserID); err != nil {
			return err
		}

		if err := h.repo.Save(txCtx, session); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "session",
			Action:    "revoked",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		return h.publisher.Publish(txCtx, session.PullEvents(), meta.ToMetadata())
	})
}
