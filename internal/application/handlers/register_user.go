package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type CreateUserHandler struct {
	uow             appContracts.UnitOfWork
	userRepo        repository.UserRepository
	passwordHasher  contracts.PasswordHasher
	passwordService contracts.PasswordService
	publisher       appContracts.EventPublisher
	sessionService  appContracts.SessionService
}

func NewCreateUserHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	passwordService contracts.PasswordService,
	publisher appContracts.EventPublisher,
	sessionService appContracts.SessionService,
) *CreateUserHandler {
	return &CreateUserHandler{
		uow:             uow,
		userRepo:        userRepo,
		passwordHasher:  passwordHasher,
		passwordService: passwordService,
		publisher:       publisher,
		sessionService:  sessionService,
	}
}

func (h *CreateUserHandler) Handle(
	ctx context.Context,
	cmd command.CreateUserCommand,
) (*dto.RegisterUserResponseDTO, error) {

	var (
		userID    string
		emailStr  string
		firstName string
		lastName  string
	)

	if err := h.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		email, err := valueobjects.NewEmail(cmd.Email)
		if err != nil {
			return err
		}

		exists, err := h.userRepo.ExistsByEmail(txCtx, email)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("email already registered")
		}

		if err := h.passwordService.Validate(cmd.Password); err != nil {
			return err
		}

		hash, err := h.passwordHasher.Hash(cmd.Password)
		if err != nil {
			return err
		}

		passwordVO, err := valueobjects.NewHashedPassword(hash)
		if err != nil {
			return err
		}

		now := time.Now().UTC()

		agg, err := aggregates.NewUserAggregate(
			email,
			passwordVO,
			cmd.FirstName,
			cmd.LastName,
			*cmd.MiddleName,
			now,
		)
		if err != nil {
			return err
		}

		if err := h.userRepo.Create(txCtx, agg); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "user",
			Action:    "created",
		}

		if err := h.publisher.Publish(
			txCtx,
			agg.PullEvents(),
			meta.ToMetadata(),
		); err != nil {
			return err
		}

		agg.ClearEvents()

		u := agg.User

		userID = u.ID()
		emailStr = u.Email().String()
		firstName = u.FirstName()
		lastName = u.LastName()

		return nil
	}); err != nil {
		return nil, err
	}

	// 2) Create session + tokens (not necessarily in the same DB transaction)
	//
	// Ideally, IP / UserAgent / DeviceID come from the outer layer (HTTP),
	// either via CreateUserCommand fields or a richer context.
	// For now, we call with empty strings as placeholders.
	sessionResult, err := h.sessionService.Create(
		ctx,
		userID,
		"", // ipAddress  (TODO: supply from transport layer)
		"", // userAgent  (TODO: supply from transport layer)
		"", // deviceID   (TODO: supply from transport layer)
	)
	if err != nil {
		// At this point the user is created but session failed.
		// You can either return 500 or handle specially.
		return nil, err
	}

	resp := &dto.RegisterUserResponseDTO{
		UserID:                userID,
		Email:                 emailStr,
		FirstName:             firstName,
		LastName:              lastName,
		AccessToken:           sessionResult.AccessToken.Value,
		AccessTokenExpiresAt:  sessionResult.AccessToken.ExpiresAt,
		RefreshToken:          sessionResult.RefreshToken.Value,
		RefreshTokenExpiresAt: sessionResult.RefreshToken.ExpiresAt,
	}

	return resp, nil
}