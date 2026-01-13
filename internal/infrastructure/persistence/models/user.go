package models

import "time"

type UserModel struct {
	ID                  string     `db:"id"`
	Email               string     `db:"email"`
	PasswordHash        string     `db:"password_hash"`
	Role                string     `db:"role"`
	FirstName           string     `db:"first_name"`
	LastName            string     `db:"last_name"`
	MiddleName          string     `db:"middle_name"`
	IsActive            bool       `db:"is_active"`
	IsVerified          bool       `db:"is_verified"`
	LastLoginAt         *time.Time `db:"last_login_at"`
	FailedLoginAttempts int        `db:"failed_login_attempts"`
	LastFailedAt        *time.Time `db:"last_failed_at"`
	LockedAt            *time.Time `db:"locked_at"`
	DeletedAt           *time.Time `db:"deleted_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
	Version             int        `db:"version"`
}
