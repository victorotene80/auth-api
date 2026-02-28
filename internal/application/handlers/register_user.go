// internal/application/handlers/create_user_handler.go
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
	roleRepo        repository.RoleRepository
	publisher       appContracts.EventPublisher
	sessionService  appContracts.SessionService
}

func NewCreateUserHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	passwordService contracts.PasswordService,
	roleRepo repository.RoleRepository,
	publisher appContracts.EventPublisher,
	sessionService appContracts.SessionService,
) *CreateUserHandler {
	return &CreateUserHandler{
		uow:             uow,
		userRepo:        userRepo,
		passwordHasher:  passwordHasher,
		passwordService: passwordService,
		roleRepo:        roleRepo,
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

		var middleName string
		if cmd.MiddleName != nil {
			middleName = *cmd.MiddleName
		}

		agg, err := aggregates.NewUserAggregate(
			email,
			passwordVO,
			cmd.FirstName,
			cmd.LastName,
			middleName,
			now,
		)
		if err != nil {
			return err
		}

		// 1) Persist user
		if err := h.userRepo.Create(txCtx, agg); err != nil {
			return err
		}

		// 2) Extract fields for later use
		u := agg.User
		userID = u.ID()
		emailStr = u.Email().String()
		firstName = u.FirstName()
		lastName = u.LastName()

		// 3) Assign default role "member"
		role, err := h.roleRepo.FindBySlug(txCtx, "member")
		if err != nil {
			return err
		}
		if role == nil {
			return fmt.Errorf("default role 'member' not found")
		}

		if err := h.roleRepo.AssignRole(
			txCtx,
			userID,
			role.ID,
			nil, // organizationID (global role)
			nil, // grantedBy
			nil, // expiresAt
		); err != nil {
			return err
		}

		// 4) Publish domain events
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
		return nil
	}); err != nil {
		return nil, err
	}

	// 5) Create session + tokens (outside DB tx)
	sessionResult, err := h.sessionService.Create(
		ctx,
		userID,
		cmd.IPAddress,
		cmd.UserAgent,
		cmd.DeviceID,
		cmd.DeviceFingerprint,
		cmd.DeviceName,
	)
	if err != nil {
		// User is created and has role, but session creation failed.
		// For now we just bubble up error (caller can decide HTTP 500 etc.).
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