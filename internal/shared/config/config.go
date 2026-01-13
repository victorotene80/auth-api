package config

/*import (
	"fmt"
	"time"
)*/
import "time"

type Config struct {
	Security SecurityConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type SecurityConfig struct {
	SessionPepper   string
	JWTSecret       string        // secret for signing JWT
	AccessTokenTTL  time.Duration // e.g., 15 * time.Minute
	RefreshTokenTTL time.Duration // e.g., 7 * 24 * time.Hour
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
