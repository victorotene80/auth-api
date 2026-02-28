package handlers

/*import (
	"context"
	"fmt"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type ChangePasswordHandler struct {
	uow            appContracts.UnitOfWork
	userRepo       repository.UserRepository
	passwordHasher contracts.PasswordHasher
	passwordSvc    services.PasswordService
	publisher      appContracts.EventPublisher
	clock          func() time.Time
}

func NewChangePasswordHandler(
	uow appContracts.UnitOfWork,
	userRepo repository.UserRepository,
	passwordHasher contracts.PasswordHasher,
	passwordSvc services.PasswordService,
	publisher appContracts.EventPublisher,
	clock func() time.Time,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		uow:            uow,
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		passwordSvc:    passwordSvc,
		publisher:      publisher,
		clock:          clock,
	}
}

func (h *ChangePasswordHandler) Handle(
	ctx context.Context,
	cmd command.ChangePasswordCommand,
) (*dto.ChangePasswordDTO, error) {

	var result *dto.ChangePasswordDTO

	err := h.uow.WithinTransaction(ctx, func(ctx context.Context) error {
		userAgg, err := h.userRepo.FindByID(ctx, cmd.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		if !h.passwordHasher.Verify(cmd.OldPassword, userAgg.User.Password().Value()) {
			return fmt.Errorf("old password does not match")
		}

		if err := h.passwordSvc.Validate(cmd.NewPassword); err != nil {
			return err
		}

		newHash, err := h.passwordHasher.Hash(cmd.NewPassword)
		if err != nil {
			return err
		}

		newHashed, err := valueobjects.NewHashedPassword(newHash)
		if err != nil {
			return err
		}

		if err := userAgg.ChangePassword(userAgg.User.Password().Value(), newHashed); err != nil {
			return err
		}

		if err := h.userRepo.Update(ctx, userAgg); err != nil {
			return err
		}

		meta := messaging.Context{
			Aggregate: "user",
			Action:    "password_changed",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}

		if err := h.publisher.Publish(ctx, userAgg.PullEvents(), meta.ToMetadata()); err != nil {
			return err
		}

		result = &dto.ChangePasswordDTO{
			UserID:    userAgg.User.ID(),
			Success:   true,
			ChangedAt: h.clock(),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}*/
