package handlers

/*import (
	"context"
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type LoginHandler struct {
	userRepo       repository.UserRepository
	sessionSvc     contracts.SessionService
	accountPolicy  policy.AccountLockPolicy
	mfaService     *services.MFAService
	eventPublisher appContracts.EventPublisher
	passwordHasher contracts.PasswordHasher
	clock          func() time.Time
	uow            appContracts.UnitOfWork
}

func NewLoginHandler(
	userRepo repository.UserRepository,
	sessionSvc contracts.SessionService,
	accountPolicy policy.AccountLockPolicy,
	mfaService *services.MFAService,
	eventPublisher appContracts.EventPublisher,
	passwordHasher contracts.PasswordHasher,
	clock func() time.Time,
	uow appContracts.UnitOfWork,
) *LoginHandler {
	return &LoginHandler{
		userRepo:       userRepo,
		sessionSvc:     sessionSvc,
		accountPolicy:  accountPolicy,
		mfaService:     mfaService,
		eventPublisher: eventPublisher,
		passwordHasher: passwordHasher,
		clock:          clock,
		uow:            uow,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd command.LoginCommand) (*dto.LoginResultDTO, error) {
	now := h.clock()

	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	var result *dto.LoginResultDTO

	err = h.uow.WithinTransaction(ctx, func(exec context.Context) error {
		userAgg, err := h.userRepo.FindByEmail(exec, email)
		if err != nil {
			return errors.New("invalid credentials")
		}
		user := userAgg.User

		if user.LockedAt != nil && h.accountPolicy.IsAccountStillLocked(*user.LockedAt) {
			return errors.New("account is locked")
		}

		if !h.passwordHasher.Verify(cmd.Password, user.Password().Value()) {
			userAgg.RecordFailedLogin(now)
			if err := h.userRepo.Update(exec, userAgg); err != nil {
				return err
			}
			return errors.New("invalid credentials")
		}

		userAgg.ResetFailedLogins()

		mfaRequired := false
		var challengeID *string
		var sessionResult contracts.SessionResult

		if h.mfaService != nil && h.mfaService.IsMFARequiredForAction("login") {
			mfaRequired = true
			challenge, err := h.mfaService.StartMFAChallenge()
			if err != nil {
				return errors.New("failed to start MFA challenge")
			}
			challengeID = &challenge
		}

		if !mfaRequired {
			sessionResult, err = h.sessionSvc.Create(
				exec,
				user.ID(),
				user.Role(),
				cmd.IPAddress,
				cmd.UserAgent,
				cmd.DeviceID,
			)
			if err != nil {
				return err
			}

			userAgg.RecordLogin(now)
			userAgg.User.SetLastLogin(now)

			if err := h.userRepo.Update(exec, userAgg); err != nil {
				return err
			}
		}

		meta := messaging.Context{
			Aggregate: "user",
			Action:    "login",
			IPAddress: cmd.IPAddress,
			UserAgent: cmd.UserAgent,
			DeviceID:  cmd.DeviceID,
		}
		if err := h.eventPublisher.Publish(exec, userAgg.PullEvents(), meta.ToMetadata()); err != nil {
			return err
		}

		status := "SUCCESS"
		if mfaRequired {
			status = "MFA_REQUIRED"
		}

		result = &dto.LoginResultDTO{
			Status: status,
			Tokens: contracts.TokenPair{
				AccessToken:  sessionResult.AccessToken,
				RefreshToken: sessionResult.RefreshToken,
			},
			MFARequired: mfaRequired,
			LastLogin:   user.LastLoginAt(),
			ChallengeID: challengeID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}*/
