package bootstrap

import (
	"time"

	"github.com/victorotene80/authentication_api/internal/application/command"
	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/application/handlers"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	appServices "github.com/victorotene80/authentication_api/internal/application/services"
	"github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	infraService "github.com/victorotene80/authentication_api/internal/infrastructure/services"
	"github.com/victorotene80/authentication_api/internal/shared/config"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
	"golang.org/x/crypto/bcrypt"
)

func initializeCommands(
	p *Persistence,
	eventPublisher appContracts.EventPublisher,
	cfg *config.Config,
) *messaging.CommandBus {

	bus := messaging.NewCommandBus()
	clock := time.Now

	passwordHasher, err := infraService.NewBcryptPasswordHasher(
		cfg.Security.SessionPepper,
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}

	sessionHasher, err := utils.NewSessionKeyHasher(cfg.Security.SessionPepper)
	if err != nil {
		panic(err)
	}

	tokenGen := infraService.NewJWTGenerator(cfg.Security.JWTSecret)
	accountPolicy := policy.DefaultAccountLockPolicy()
	sessionPolicy := policy.DefaultSessionPolicy()
	passwordService := services.NewPasswordService(policy.DefaultPasswordPolicy())
	mfaService := services.NewMFAService(policy.DefaultMFAPolicy())

	sessionSvc := appServices.NewSessionService(
		p.SessionRepo,
		tokenGen,
		p.UoW,
		sessionHasher,
		sessionPolicy,
		clock,
		eventPublisher,
	)

	messaging.MustRegister[command.LoginCommand, *dto.LoginResultDTO](
		bus,
		handlers.NewLoginHandler(
			p.UserRepo,
			sessionSvc,
			accountPolicy,
			mfaService,
			eventPublisher,
			passwordHasher,
			clock,
			p.UoW,
		),
	)

	messaging.MustRegister[command.CreateUserCommand, *dto.UserResponseDTO](
		bus,
		handlers.NewCreateUserHandler(
			p.UoW,
			p.UserRepo,
			passwordHasher,
			passwordService,
			eventPublisher,
		),
	)

	messaging.MustRegister[command.LogoutCommand, *dto.LogoutDTO](
		bus,
		handlers.NewLogoutHandler(
			p.UoW,
			p.SessionRepo,
			eventPublisher,
			clock,
		),
	)

	return bus
}
