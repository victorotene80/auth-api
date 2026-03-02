package models

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID                    string         `db:"id"`
	Email                 string         `db:"email"`
	PasswordHash          string         `db:"password_hash"`
	FirstName             string         `db:"first_name"`
	LastName              string         `db:"last_name"`
	MiddleName            sql.NullString `db:"middle_name"`
	Status                string         `db:"status"`
	EmailVerified         bool           `db:"email_verified"`
	EmailVerifiedAt       sql.NullTime   `db:"email_verified_at"`
	PasswordChangedAt     sql.NullTime   `db:"password_changed_at"`
	PasswordExpiresAt     sql.NullTime   `db:"password_expires_at"`
	RequirePasswordChange bool           `db:"require_password_change"`
	FailedLoginAttempts   int            `db:"failed_login_attempts"`
	LockedUntil           sql.NullTime   `db:"locked_until"`
	LastLoginAt           sql.NullTime   `db:"last_login_at"`
	LastLoginIP           sql.NullString `db:"last_login_ip"`
	LastActiveAt          sql.NullTime   `db:"last_active_at"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
	DeletedAt             sql.NullTime   `db:"deleted_at"`
}