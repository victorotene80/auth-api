package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	redisCfg, err := loadRedis()
	if err != nil {
		return nil, err
	}

	geo, err := loadGeoIP()
	if err != nil {
		return nil, err
	}

	messaging := loadMessaging()

	return &Config{
		Security:  security,
		Database:  database,
		Redis:     redisCfg,
		GeoIP:     geo,
		Messaging: messaging,
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

	return SecurityConfig{
		SessionPepper:   pepper,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  getDurationOrDefault("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getDurationOrDefault("REFRESH_TOKEN_TTL", 7*24*time.Hour),
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

	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: getStringOrDefault("REDIS_PASSWORD", ""),
		DB:       getIntOrDefault("REDIS_DB", 0),
		TTL:      getDurationOrDefault("REDIS_TTL", 15*time.Minute),
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
		MaxLifetime:  getIntOrDefault("DB_MAX_LIFETIME", 300),
	}, nil
}

func loadGeoIP() (GeoIPConfig, error) {
	path, err := getString("GEOIP_DB_PATH")
	if err != nil {
		return GeoIPConfig{}, err
	}

	return GeoIPConfig{DBPath: path}, nil
}

func loadMessaging() MessagingConfig {
	return MessagingConfig{
		Kafka: KafkaConfig{
			Brokers:         splitCSV(getStringOrDefault("KAFKA_BROKERS", "localhost:9092")),
			TopicPrefix:     getStringOrDefault("KAFKA_TOPIC_PREFIX", "auth."),
			ConsumerGroupID: getStringOrDefault("KAFKA_CONSUMER_GROUP_ID", "auth-service"),
			WriteTimeout:    getDurationOrDefault("KAFKA_WRITE_TIMEOUT", 10*time.Second),
		},
		RabbitMQ: RabbitMQConfig{
			DSN:            getStringOrDefault("RABBITMQ_DSN", "amqp://guest:guest@localhost:5672/"),
			Exchange:       getStringOrDefault("RABBITMQ_EXCHANGE", "auth.tasks"),
			RetryExchange:  getStringOrDefault("RABBITMQ_RETRY_EXCHANGE", "auth.tasks.retry"),
			DLExchange:     getStringOrDefault("RABBITMQ_DL_EXCHANGE", "auth.tasks.dlx"),
			MaxRetries:     getIntOrDefault("RABBITMQ_MAX_RETRIES", 3),
			RetryDelay:     getDurationOrDefault("RABBITMQ_RETRY_DELAY", 15*time.Second),
			PublishTimeout: getDurationOrDefault("RABBITMQ_PUBLISH_TIMEOUT", 5*time.Second),
		},
		Relay: RelayConfig{
			PollInterval:       getDurationOrDefault("RELAY_POLL_INTERVAL", 1*time.Second),
			BatchSize:          getIntOrDefault("RELAY_BATCH_SIZE", 50),
			ReclaimAfter:       getDurationOrDefault("RELAY_RECLAIM_AFTER", 2*time.Minute),
			DefaultEventBroker: "event",
			DefaultTaskBroker:  "task",
			EventRoutes: map[string]string{
				"auth.user.created.v1":     "event",
				"auth.user.locked.v1":      "event",
				"auth.session.created.v1":  "event",
				"auth.session.revoked.v1":  "event",
				"auth.password.changed.v1": "event",
			},
			TaskRoutes: map[string]string{
				"auth.send-welcome-email.v1":  "task",
				"auth.send-verification.v1":   "task",
				"auth.sync-analytics-user.v1": "task",
			},
		},
	}
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
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