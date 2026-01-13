package handlers

import (
	"context"
	"fmt"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	//"github.com/victorotene80/authentication_api/internal/application/messaging"
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
}

func NewCreateUserHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	passwordService contracts.PasswordService,
	publisher appContracts.EventPublisher,
) *CreateUserHandler {
	return &CreateUserHandler{
		uow:             uow,
		userRepo:        userRepo,
		passwordHasher:  passwordHasher,
		passwordService: passwordService,
		publisher:       publisher,
	}
}

func (h *CreateUserHandler) Handle(
	ctx context.Context,
	cmd command.CreateUserCommand,
) (*dto.UserResponseDTO, error) {

	var result *dto.UserResponseDTO

	err := h.uow.WithinTransaction(ctx, func(ctx context.Context) error {
		email, err := valueobjects.NewEmail(cmd.Email)
		if err != nil {
			return err
		}

		exists, err := h.userRepo.ExistsByEmail(ctx, email)
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

		password, err := valueobjects.NewHashedPassword(hash)
		if err != nil {
			return err
		}

		role, err := valueobjects.NewRole("user")
		if err != nil {
			return err
		}

		agg, err := aggregates.NewUserAggregate(
			email,
			password,
			cmd.FirstName,
			cmd.LastName,
			cmd.MiddleName,
			role,
		)
		if err != nil {
			return err
		}

		if err := h.userRepo.Create(ctx, agg); err != nil {
			return err
		}

		//msgs := make([]messaging.EventMessage, 0)

		//for _, e := range agg.PullEvents(){
		//	msgs = append(msgs, messaging.ToEventMessage(e))
		//}

		meta := messaging.Context{
			Aggregate: "user",
			Action: "created",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(
			ctx,
			agg.PullEvents(),
			meta.ToMetadata(),
		); err != nil {
			return err
		}

		u := agg.User
		result = &dto.UserResponseDTO{
			UserID:    u.ID(),
			Email:     u.Email().String(),
			FirstName: u.FirstName(),
			LastName:  u.LastName(),
		}

		return nil
	})

	if err != nil {
		return &dto.UserResponseDTO{}, err
	}

	return result, nil
}
