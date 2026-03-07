package config

import "time"

type Config struct {
	Security  SecurityConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	GeoIP     GeoIPConfig
	Messaging MessagingConfig
}

type MessagingConfig struct {
	Kafka    KafkaConfig
	RabbitMQ RabbitMQConfig
	Relay    RelayConfig
}

type KafkaConfig struct {
	Brokers         []string
	TopicPrefix     string
	ConsumerGroupID string
	WriteTimeout    time.Duration
}

type RabbitMQConfig struct {
	DSN            string
	Exchange       string
	RetryExchange  string
	DLExchange     string
	MaxRetries     int
	RetryDelay     time.Duration
	PublishTimeout time.Duration
}

type RelayConfig struct {
	PollInterval        time.Duration
	BatchSize           int
	ReclaimAfter        time.Duration
	DefaultEventBroker  string
	DefaultTaskBroker   string
	EventRoutes         map[string]string
	TaskRoutes          map[string]string
}

type SecurityConfig struct {
	SessionPepper   string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	TTL      time.Duration
}

type GeoIPConfig struct {
	DBPath string
}