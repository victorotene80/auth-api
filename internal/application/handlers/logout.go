package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
)

type LogoutHandler struct {
	uow         appContracts.UnitOfWork
	sessionRepo repository.SessionRepository
	publisher   appContracts.EventPublisher
	clock       func() time.Time
}

func NewLogoutHandler(
	uow appContracts.UnitOfWork,
	sessionRepo repository.SessionRepository,
	publisher appContracts.EventPublisher,
	clock func() time.Time,
) *LogoutHandler {
	return &LogoutHandler{
		uow:         uow,
		sessionRepo: sessionRepo,
		publisher:   publisher,
		clock:       clock,
	}
}

func (h *LogoutHandler) Handle(ctx context.Context, cmd command.LogoutCommand) (*dto.LogoutDTO, error) {
	var result *dto.LogoutDTO
	err := h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		now := h.clock()
		session, err := h.sessionRepo.FindByID(txCtx, cmd.SessionID)
		if err != nil {
			return fmt.Errorf("session not found: %w", err)
		}

		if session.UserID() != cmd.UserID {
			return fmt.Errorf("session does not belong to user")
		}

		if !session.IsRevoked() {
			reason := cmd.Reason
			if reason == "" {
				reason = "user logout"
			}
			if err := session.Revoke(now, reason); err != nil {
				return fmt.Errorf("cannot revoke session: %w", err)
			}
			if err := h.sessionRepo.Save(txCtx, session); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}
		}

		meta := messaging.Context{
			Aggregate: "session",
			Action:    "revoked",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(txCtx, session.PullEvents(), meta.ToMetadata()); err != nil {
			return fmt.Errorf("failed to publish events: %w", err)
		}

		session.ClearEvents()

		result = &dto.LogoutDTO{
			SessionID: session.ID(),
			Status:    "SUCCESS",
			Reason:    cmd.Reason,
			Time:      now,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
