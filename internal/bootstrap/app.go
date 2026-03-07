package bootstrap

import (
	"context"
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
	Stop   func()
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
		return nil, err
	}

	redisClient, err := initializeRedis(cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	messagePublisher := initializeMessagePublisher(persistenceLayer, logger)

	consumers, stopMessaging := initializeMessaging(persistenceLayer, cfg, logger)
	registerConsumers(consumers, logger)

	consumeCtx, cancelConsume := context.WithCancel(context.Background())

	go func() {
		if err := consumers.EventConsumer.Consume(consumeCtx); err != nil {
			logger.Error("event consumer exited", zap.Error(err))
		}
	}()
	go func() {
		if err := consumers.TaskConsumer.Consume(consumeCtx); err != nil {
			logger.Error("task consumer exited", zap.Error(err))
		}
	}()

	commandBus, authSvc := initializeCommands(
		persistenceLayer,
		messagePublisher,
		cfg,
		redisClient,
		logger,
	)

	router := initializeHTTP(commandBus, authSvc, logger)

	stop := func() {
		cancelConsume()
		stopMessaging()
		logger.Info("app shutdown complete")
	}

	return &App{
		Router: router,
		DB:     persistenceLayer.DB,
		Redis:  redisClient,
		Stop:   stop,
	}, nil
}

/*package bootstrap

import (
	"context"
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
	Stop func()
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
		return nil, err
	}

	redisClient, err := initializeRedis(cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	eventPublisher := initializeEventPublisher(persistenceLayer, logger)

	consumers, stopMessaging := initializeMessaging(persistenceLayer, cfg, logger)

	// All topic → handler wiring lives in consumers.go, not here.
	registerConsumers(consumers, logger)

	// Start consumer loops — each blocks, so they run in goroutines.
	consumeCtx, cancelConsume := context.WithCancel(context.Background())

	go func() {
		if err := consumers.EventConsumer.Consume(consumeCtx); err != nil {
			logger.Error("event consumer exited", zap.Error(err))
		}
	}()

	go func() {
		if err := consumers.TaskConsumer.Consume(consumeCtx); err != nil {
			logger.Error("task consumer exited", zap.Error(err))
		}
	}()

	commandBus, authSvc := initializeCommands(
		persistenceLayer,
		eventPublisher,
		cfg,
		redisClient,
		logger,
	)

	router := initializeHTTP(commandBus, authSvc, logger)

	stop := func() {
		cancelConsume()
		stopMessaging()
		logger.Info("app shutdown complete")
	}

	return &App{
		Router: router,
		DB:     persistenceLayer.DB,
		Redis:  redisClient,
		Stop:   stop,
	}, nil
}*/