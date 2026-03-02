// internal/bootstrap/app.go
package bootstrap

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/victorotene80/authentication_api/internal/shared/config"
	"github.com/victorotene80/authentication_api/internal/shared/logging"
)

type App struct {
	Router http.Handler
	DB     any
	Redis  *redis.Client
}

func InitializeApp() (*App, error) {
	logProvider, err := logging.NewLoggerProvider()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logProvider.Sync()

	logger := logProvider.Logger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
		return nil, err
	}

	persistenceLayer, err := initializePersistence(cfg, logger)
	if err != nil {
		//logger.Fatal("persistence initialization failed", zap.Error(err))
		return nil, err
	}

	redisClient, err := initializeRedis(cfg.Redis, logger)
	if err != nil {
		//logger.Fatal("redis initialization failed", zap.Error(err))
		return nil, err
	}

	eventPublisher := initializeEventPublisher(persistenceLayer, logger)

	commandBus, authSvc := initializeCommands(
		persistenceLayer,
		eventPublisher,
		cfg,
		redisClient,
		logger,
	)

	router := initializeHTTP(commandBus, authSvc, logger)

	return &App{
		Router: router,
		DB:     persistenceLayer.DB,
		Redis:  redisClient,
	}, nil
}