package bootstrap

import (
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	appHandlers "github.com/victorotene80/authentication_api/internal/application/handlers"
	"github.com/victorotene80/authentication_api/internal/application/messaging"
	msgmw "github.com/victorotene80/authentication_api/internal/application/messaging/middleware"
	appServices "github.com/victorotene80/authentication_api/internal/application/services"

	domainServices "github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"

	audit "github.com/victorotene80/authentication_api/internal/infrastructure/persistence"
	cacheInfra "github.com/victorotene80/authentication_api/internal/infrastructure/persistence/cache"
	infraServices "github.com/victorotene80/authentication_api/internal/infrastructure/services"

	"github.com/victorotene80/authentication_api/internal/shared/config"
	"github.com/victorotene80/authentication_api/internal/shared/utils"

	"golang.org/x/crypto/bcrypt"
)

func initializeCommands(
	p *Persistence,
	eventPublisher appContracts.EventPublisher,
	cfg *config.Config,
	redisClient *redis.Client,
	logger *zap.Logger,
) (*messaging.CommandBus, appContracts.AuthService) {

	bus := messaging.NewCommandBus()
	clock := func() time.Time { return time.Now().UTC() }

	if logger != nil {
		msgmw.AttachLogging(bus, logger)
	}

	passwordHasher, err := infraServices.NewBcryptPasswordHasher(
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

	tokenGen := infraServices.NewJWTGenerator(cfg.Security.JWTSecret)

	sessionPolicy := policy.DefaultSessionPolicy()
	passwordPolicy := policy.DefaultPasswordPolicy()
	passwordService := domainServices.NewPasswordService(passwordPolicy)
	roleRepo := p.RoleRepo

		lockPolicy := policy.DefaultAccountLockPolicy()
	lockService := domainServices.NewAccountLockService(lockPolicy)

	geoIPService, err := infraServices.NewMaxMindGeoIPService(cfg.GeoIP.DBPath)
	if err != nil {
		panic(err)
	}

	auditLogger := audit.NewPostgresAuditLogger(p.DB)

	sessionCache := cacheInfra.NewRedisCache[string, cacheInfra.CachedSession](
		redisClient,
		"session:",
		sessionPolicy.MaxDuration,
	)

	sessionSvc := appServices.NewSessionService(
		p.SessionRepo,
		tokenGen,
		p.UoW,
		sessionHasher,
		sessionPolicy,
		clock,
		eventPublisher,
		geoIPService,
		sessionCache,
		auditLogger,
	)

	authSessionCache := cacheInfra.NewRedisCache[string, cacheInfra.CachedSession](
		redisClient,
		"auth-session:",
		sessionPolicy.MaxDuration,
	)

	authSvc := appServices.NewAuthService(
		tokenGen,
		p.SessionRepo,
		sessionPolicy,
		authSessionCache,
		auditLogger,
	)

	createUserHandler := appHandlers.NewCreateUserHandler(
		p.UoW,
		p.UserRepo,
		passwordHasher,
		passwordService,
		roleRepo,
		eventPublisher,
		sessionSvc,
	)

	loginHandler := appHandlers.NewLoginHandler(
		p.UoW,
		p.UserRepo,
		passwordHasher,
		sessionSvc,
		lockService,
		eventPublisher,
		clock,
	)

	messaging.MustRegister(bus, loginHandler)
	messaging.MustRegister(bus, createUserHandler)

	return bus, authSvc
}