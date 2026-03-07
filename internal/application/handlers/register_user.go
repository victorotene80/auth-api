package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	publisher       appContracts.MessagePublisher
	sessionService  appContracts.SessionService
}

func NewCreateUserHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	passwordService contracts.PasswordService,
	roleRepo repository.RoleRepository,
	publisher appContracts.MessagePublisher,
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

		var phoneNumber valueobjects.PhoneNumber
		if cmd.Phone != nil && *cmd.Phone != "" {
			phoneNumber, err = valueobjects.NewPhoneNumber(*cmd.Phone)
			if err != nil {
				return err
			}
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
			cmd.IPAddress,
			phoneNumber,
			now,
		)
		if err != nil {
			return err
		}

		if err := h.userRepo.Create(txCtx, agg); err != nil {
			return err
		}

		u := agg.User
		userID = u.ID()
		emailStr = u.Email().String()
		firstName = u.FirstName()
		lastName = u.LastName()

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
			nil,
			nil,
			nil,
		); err != nil {
			return err
		}

		userCreatedMeta := messaging.Context{
			Kind:          messaging.KindIntegrationEvent,
			Name:          "auth.user.created.v1",
			AggregateType: "user",
			Action:        "created",
		}

		if err := h.publisher.Publish(
			txCtx,
			agg.PullEvents(),
			userCreatedMeta.ToMetadata(),
		); err != nil {
			return err
		}

		roleAssignedEnv := messaging.Envelope{
			ID:            uuid.NewString(),
			Name:          "auth.user.role-assigned.v1",
			Kind:          messaging.KindIntegrationEvent,
			AggregateID:   userID,
			AggregateType: "user",
			OccurredAt:    now,
			Payload: messaging.MustJSON(map[string]any{
				"user_id":   userID,
				"role_id":   role.ID,
				"role_slug": role.Slug,
				"role_name": role.Name,
				"version":   1,
			}),
			Metadata: map[string]string{
				"message_kind":   string(messaging.KindIntegrationEvent),
				"message_name":   "auth.user.role-assigned.v1",
				"aggregate_type": "user",
				"action":         "role_assigned",
			},
			Version: 1,
		}

		if err := h.publisher.PublishEnvelope(txCtx, roleAssignedEnv); err != nil {
			return err
		}

		taskEnv := messaging.Envelope{
			ID:            uuid.NewString(),
			Name:          "auth.send-welcome-email.v1",
			Kind:          messaging.KindTask,
			AggregateID:   userID,
			AggregateType: "user",
			OccurredAt:    now,
			Payload: messaging.MustJSON(map[string]any{
				"user_id": userID,
				"email":   emailStr,
			}),
			Metadata: map[string]string{
				"message_kind":   string(messaging.KindTask),
				"message_name":   "auth.send-welcome-email.v1",
				"aggregate_type": "user",
			},
			Version: 1,
		}

		if err := h.publisher.PublishEnvelope(txCtx, taskEnv); err != nil {
			return err
		}

		agg.ClearEvents()
		return nil
	}); err != nil {
		return nil, err
	}

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
