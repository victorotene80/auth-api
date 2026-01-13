package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/victorotene80/authentication_api/internal/shared/config"
)

func initializeRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}

	return client, nil
}
