package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type LoginHandler struct {
	uow            appContracts.UnitOfWork
	userRepo       repository.UserRepository
	passwordHasher contracts.PasswordHasher
	sessionService appContracts.SessionService
	lockService    *services.AccountLockService
	eventPublisher appContracts.MessagePublisher
	clock          func() time.Time
}

func NewLoginHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	sessionService appContracts.SessionService,
	lockService *services.AccountLockService,
	eventPublisher appContracts.MessagePublisher,
	clock func() time.Time,
) *LoginHandler {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}

	return &LoginHandler{
		uow:            uow,
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		sessionService: sessionService,
		lockService:    lockService,
		eventPublisher: eventPublisher,
		clock:          clock,
	}
}

func (h *LoginHandler) Handle(
	ctx context.Context,
	cmd command.LoginCommand,
) (*dto.LoginResultDTO, error) {

	now := h.clock()

	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	var (
		userAgg   *aggregates.UserAggregate
		lastLogin *time.Time
	)

	if err := h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		agg, err := h.userRepo.FindByEmail(txCtx, email)
		if err != nil {
			return errors.New("invalid credentials")
		}

		userAgg = agg
		user := agg.User

		if !userAgg.EnsureNotLocked(now, h.lockService) {
			return errors.New("account is locked")
		}

		if !h.passwordHasher.Verify(cmd.Password, user.Password().Value()) {
			userAgg.RecordFailedLogin(now, h.lockService)

			if err := h.userRepo.Update(txCtx, userAgg); err != nil {
				return err
			}

			return errors.New("invalid credentials")
		}

		userAgg.RecordLogin(now, cmd.IPAddress)
		lastLogin = user.LastLoginAt()

		if err := h.userRepo.Update(txCtx, userAgg); err != nil {
			return err
		}

		/*meta := messaging.Context{
			Aggregate: "user",
			Action:    "login",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}*/

		meta := messaging.Context{
			Kind:          messaging.KindIntegrationEvent,
			Name:          "auth.user.login-recorded.v1",
			AggregateType: "user",
			Action:        "login",
			IPAddress:     cmd.IPAddress,
			UserAgent:     cmd.UserAgent,
			DeviceID:      cmd.DeviceID,
		}

		if err := h.eventPublisher.Publish(
			txCtx,
			userAgg.PullEvents(),
			meta.ToMetadata(),
		); err != nil {
			return err
		}

		userAgg.ClearEvents()
		return nil
	}); err != nil {
		return nil, err
	}

	sessionResult, err := h.sessionService.Create(
		ctx,
		userAgg.User.ID(),
		cmd.IPAddress,
		cmd.UserAgent,
		cmd.DeviceID,
		cmd.DeviceFingerprint,
		cmd.DeviceName,
	)
	if err != nil {
		return nil, err
	}

	result := &dto.LoginResultDTO{
		Status:      "SUCCESS",
		MFARequired: false,     // wire MFA later if you want
		LastLogin:   lastLogin, // may be nil if first login
		ChallengeID: nil,       // for MFA
		Tokens: contracts.TokenPair{
			AccessToken: contracts.Token{
				Value:     sessionResult.AccessToken.Value,
				ExpiresAt: sessionResult.AccessToken.ExpiresAt,
			},
			RefreshToken: contracts.Token{
				Value:     sessionResult.RefreshToken.Value,
				ExpiresAt: sessionResult.RefreshToken.ExpiresAt,
			},
		},
	}

	return result, nil
}
