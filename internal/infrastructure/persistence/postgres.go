package persistence

import (
	"database/sql"
	"fmt"
	"github.com/victorotene80/authentication_api/internal/shared/config"
	_ "github.com/lib/pq"
	"time"
)

func NewDatabase(cfg config.DatabaseConfig) (*sql.DB, error){
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil{
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}