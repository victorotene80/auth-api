package bootstrap

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/victorotene80/authentication_api/internal/shared/config"
)

type App struct {
	Router http.Handler
	DB     any
	Redis  *redis.Client
}

func InitializeApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	persistenceLayer, err := initializePersistence(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := initializeRedis(cfg.Redis)
	if err != nil {
		return nil, err
	}

	eventPublisher := initializeEventPublisher(persistenceLayer)
	commandBus := initializeCommands(persistenceLayer, eventPublisher, cfg)
	router := initializeHTTP(commandBus)

	return &App{
		Router: router,
		DB:     persistenceLayer.DB,
		Redis:  redisClient,
	}, nil
}
