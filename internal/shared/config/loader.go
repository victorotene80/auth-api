package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func Load() (*Config, error) {
	security, err := loadSecurity()
	if err != nil {
		return nil, err
	}

	database, err := loadDatabase()
	if err != nil {
		return nil, err
	}

	redis, err := loadRedis()
	if err != nil {
		return nil, err
	}

	return &Config{
		Security: security,
		Database: database,
		Redis: redis,
	}, nil
}

func loadSecurity() (SecurityConfig, error) {
	pepper, err := getString("SESSION_PEPPER")
	if err != nil {
		return SecurityConfig{}, err
	}

	jwtSecret, err := getString("JWT_SECRET")
	if err != nil {
		return SecurityConfig{}, err
	}

	accessTTL := getDurationOrDefault("ACCESS_TOKEN_TTL", 15*time.Minute)
	refreshTTL := getDurationOrDefault("REFRESH_TOKEN_TTL", 7*24*time.Hour)

	return SecurityConfig{
		SessionPepper:   pepper,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}, nil
}

func loadRedis() (RedisConfig, error) {
	host, err := getString("REDIS_HOST")
	if err != nil {
		return RedisConfig{}, err
	}

	port, err := getInt("REDIS_PORT")
	if err != nil {
		return RedisConfig{}, err
	}

	password := getStringOrDefault("REDIS_PASSWORD", "")
	db := getIntOrDefault("REDIS_DB", 0)
	ttl := getDurationOrDefault("REDIS_TTL", 15*time.Minute)

	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
		TTL:      ttl,
	}, nil
}

func loadDatabase() (DatabaseConfig, error) {
	host, err := getString("DB_HOST")
	if err != nil {
		return DatabaseConfig{}, err
	}

	port, err := getInt("DB_PORT")
	if err != nil {
		return DatabaseConfig{}, err
	}

	user, err := getString("DB_USER")
	if err != nil {
		return DatabaseConfig{}, err
	}

	password, err := getString("DB_PASSWORD")
	if err != nil {
		return DatabaseConfig{}, err
	}

	name, err := getString("DB_NAME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		Name:         name,
		SSLMode:      getStringOrDefault("DB_SSLMODE", "disable"),
		MaxOpenConns: getIntOrDefault("DB_MAX_OPEN", 10),
		MaxIdleConns: getIntOrDefault("DB_MAX_IDLE", 5),
	}, nil
}

func getDurationOrDefault(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func getString(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func getStringOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getInt(key string) (int, error) {
	v, err := getString(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}

	return i, nil
}

func getIntOrDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
