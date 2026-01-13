package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure"
)

type RedisSessionCache struct {
	client *redis.Client
	prefix string
}

func NewRedisSessionCache(client *redis.Client) *RedisSessionCache {
	return &RedisSessionCache{
		client: client,
		prefix: "session:",
	}
}

func (r *RedisSessionCache) key(tokenHash string) string {
	return r.prefix + tokenHash
}

func (r *RedisSessionCache) Set(ctx context.Context, session *aggregates.SessionAggregate) error {
	cached := mapAggregateToCached(session)

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	ttl := time.Until(session.ExpiresAt())
	if ttl <= 0 {
		_ = r.Delete(ctx, session.TokenHash().Value())
		return infrastructure.ErrSessionExpired
	}

	return r.client.Set(ctx, r.key(cached.TokenHash), data, ttl).Err()
}

func (r *RedisSessionCache) Get(ctx context.Context, tokenHash string) (*aggregates.SessionAggregate, error) {
	data, err := r.client.Get(ctx, r.key(tokenHash)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	var cached CachedSession
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return mapCachedToAggregate(cached)
}

func (r *RedisSessionCache) Delete(ctx context.Context, tokenHash string) error {
	return r.client.Del(ctx, r.key(tokenHash)).Err()
}

func toStringPtr(vo *valueobjects.SessionTokenHash) *string {
	if vo == nil {
		return nil
	}
	val := vo.Value()
	return &val
}

func mapAggregateToCached(s *aggregates.SessionAggregate) CachedSession {
	return CachedSession{
		ID:                s.ID(),
		UserID:            s.UserID(),
		TokenHash:         s.TokenHash().Value(),
		PreviousTokenHash: toStringPtr(s.PreviousTokenHash()),
		RotationID:        s.RotationID(),
		IPAddress:         s.IPAddress(),
		DeviceID:          s.DeviceID(),
		UserAgent:         s.UserAgent(),
		Status:            string(s.Status()),
		CreatedAt:         s.CreatedAt(),
		LastSeenAt:        s.LastSeenAt(),
		ExpiresAt:         s.ExpiresAt(),
		RevokedAt:         s.RevokedAt(),
		Version:           s.Version(),
	}
}

func mapCachedToAggregate(c CachedSession) (*aggregates.SessionAggregate, error) {
	return aggregates.RehydrateSession(
		c.ID,
		c.UserID,
		valueobjects.Role(c.Role),
		c.TokenHash,
		c.PreviousTokenHash,
		c.RotationID,
		c.IPAddress,
		c.UserAgent,
		c.DeviceID,
		c.Status,
		c.CreatedAt,
		c.LastSeenAt,
		c.ExpiresAt,
		c.RevokedAt,
		c.Version,
	)
}
