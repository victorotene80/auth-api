package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

var _ repository.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	agg *aggregates.UserAggregate,
) error {
	exec := ChooseExecutor(ctx, r.db)
	u := agg.User

	const q = `
		INSERT INTO auth.users (
			id,
			email,
			password_hash,
			first_name,
			last_name,
			middle_name,
			status,
			email_verified,
			email_verified_at,
			password_changed_at,
			password_expires_at,
			require_password_change,
			failed_login_attempts,
			locked_until,
			last_login_at,
			last_login_ip,
			last_active_at,
			created_at,
			updated_at,
			deleted_at,
			phone
		)
		VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13, $14, $15,
			$16, $17,
			$18, $19, $20, $21
		)
	`

	_, err := exec.ExecContext(ctx, q,
		u.ID(),
		u.Email().String(),
		u.Password().Value(),
		u.FirstName(),
		u.LastName(),
		nullIfEmpty(u.MiddleName()),
		u.Status().String(),
		u.EmailVerified(),
		u.EmailVerifiedAt(),
		u.PasswordChangedAt(),
		u.PasswordExpiresAt(),
		u.RequirePasswordChange(),
		u.FailedLoginAttempts(),
		u.LockedUntil(),
		u.LastLoginAt(),
		nullIfEmpty(u.LastLoginIP()),
		u.LastActiveAt(),
		u.CreatedAt(),
		u.UpdatedAt(),
		u.DeletedAt(),
		u.Phone().String(),
	)
	return err
}

func (r *PostgresUserRepository) Update(
	ctx context.Context,
	agg *aggregates.UserAggregate,
) error {
	exec := ChooseExecutor(ctx, r.db)
	u := agg.User

	const q = `
		UPDATE auth.users SET
			email                   = $1,
			password_hash           = $2,
			first_name              = $3,
			last_name               = $4,
			middle_name             = $5,
			status                  = $6,
			email_verified          = $7,
			email_verified_at       = $8,
			password_changed_at     = $9,
			password_expires_at     = $10,
			require_password_change = $11,
			failed_login_attempts   = $12,
			locked_until            = $13,
			last_login_at           = $14,
			last_login_ip           = $15,
			last_active_at          = $16,
			updated_at              = $17,
			phone                   = $18
		WHERE id = $19
	`

	_, err := exec.ExecContext(ctx, q,
		u.Email().String(),
		u.Password().Value(),
		u.FirstName(),
		u.LastName(),
		nullIfEmpty(u.MiddleName()),
		u.Status().String(),
		u.EmailVerified(),
		u.EmailVerifiedAt(),
		u.PasswordChangedAt(),
		u.PasswordExpiresAt(),
		u.RequirePasswordChange(),
		u.FailedLoginAttempts(),
		u.LockedUntil(),
		u.LastLoginAt(),
		nullIfEmpty(u.LastLoginIP()),
		u.LastActiveAt(),
		u.UpdatedAt(),
		u.Phone().String(),
		u.ID(),
	)
	return err
}

func (r *PostgresUserRepository) SoftDelete(
	ctx context.Context,
	id string,
) error {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
		UPDATE auth.users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := exec.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
}

func (r *PostgresUserRepository) FindByID(
	ctx context.Context,
	id string,
) (*aggregates.UserAggregate, error) {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			middle_name,
			status,
			email_verified,
			email_verified_at,
			password_changed_at,
			password_expires_at,
			require_password_change,
			failed_login_attempts,
			locked_until,
			last_login_at,
			last_login_ip,
			last_active_at,
			created_at,
			updated_at,
			deleted_at,
			phone
		FROM auth.users
		WHERE id = $1
	`

	row := exec.QueryRowContext(ctx, q, id)

	var m models.UserModel
	if err := row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.FirstName,
		&m.LastName,
		&m.MiddleName,
		&m.Status,
		&m.EmailVerified,
		&m.EmailVerifiedAt,
		&m.PasswordChangedAt,
		&m.PasswordExpiresAt,
		&m.RequirePasswordChange,
		&m.FailedLoginAttempts,
		&m.LockedUntil,
		&m.LastLoginAt,
		&m.LastLoginIP,
		&m.LastActiveAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
		&m.Phone,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userAggregateFromModel(&m)
}

func (r *PostgresUserRepository) FindByEmail(
	ctx context.Context,
	email valueobjects.Email,
) (*aggregates.UserAggregate, error) {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			middle_name,
			status,
			email_verified,
			email_verified_at,
			password_changed_at,
			password_expires_at,
			require_password_change,
			failed_login_attempts,
			locked_until,
			last_login_at,
			last_login_ip,
			last_active_at,
			created_at,
			updated_at,
			deleted_at,
			phone
		FROM auth.users
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	row := exec.QueryRowContext(ctx, q, email.String())

	var m models.UserModel
	if err := row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.FirstName,
		&m.LastName,
		&m.MiddleName,
		&m.Status,
		&m.EmailVerified,
		&m.EmailVerifiedAt,
		&m.PasswordChangedAt,
		&m.PasswordExpiresAt,
		&m.RequirePasswordChange,
		&m.FailedLoginAttempts,
		&m.LockedUntil,
		&m.LastLoginAt,
		&m.LastLoginIP,
		&m.LastActiveAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt,
		&m.Phone,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userAggregateFromModel(&m)
}

func (r *PostgresUserRepository) ExistsByEmail(
	ctx context.Context,
	email valueobjects.Email,
) (bool, error) {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
		SELECT EXISTS(
			SELECT 1 FROM auth.users
			WHERE email = $1
			  AND deleted_at IS NULL
		)
	`

	var exists bool
	if err := exec.QueryRowContext(ctx, q, email.String()).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresUserRepository) List(
	ctx context.Context,
	limit int,
	cursor *string,
	status *valueobjects.UserStatus,
) ([]*aggregates.UserAggregate, *string, error) {
	if limit <= 0 {
		limit = 20
	}

	exec := ChooseExecutor(ctx, r.db)

	args := []any{}
	arg := 1

	q := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			middle_name,
			status,
			email_verified,
			email_verified_at,
			password_changed_at,
			password_expires_at,
			require_password_change,
			failed_login_attempts,
			locked_until,
			last_login_at,
			last_login_ip,
			last_active_at,
			created_at,
			updated_at,
			deleted_at, 
			phone
		FROM auth.users
		WHERE deleted_at IS NULL
	`

	if status != nil {
		q += fmt.Sprintf(" AND status = $%d", arg)
		args = append(args, status.String())
		arg++
	}

	if cursor != nil && *cursor != "" {
		q += fmt.Sprintf(" AND id > $%d", arg)
		args = append(args, *cursor)
		arg++
	}

	q += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d", arg)
	args = append(args, limit+1) // limit+1 to detect "has next"

	rows, err := exec.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var (
		users      []*aggregates.UserAggregate
		lastUserID string
		count      int
	)

	for rows.Next() {
		var m models.UserModel
		if err := rows.Scan(
			&m.ID,
			&m.Email,
			&m.PasswordHash,
			&m.FirstName,
			&m.LastName,
			&m.MiddleName,
			&m.Status,
			&m.EmailVerified,
			&m.EmailVerifiedAt,
			&m.PasswordChangedAt,
			&m.PasswordExpiresAt,
			&m.RequirePasswordChange,
			&m.FailedLoginAttempts,
			&m.LockedUntil,
			&m.LastLoginAt,
			&m.LastLoginIP,
			&m.LastActiveAt,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
			&m.Phone,
		); err != nil {
			return nil, nil, err
		}

		count++
		if count > limit {
			lastUserID = m.ID
			break
		}

		agg, err := userAggregateFromModel(&m)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, agg)
		lastUserID = m.ID
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var nextCursor *string
	if count > limit && lastUserID != "" {
		nextCursor = &lastUserID
	}

	return users, nextCursor, nil
}

func userAggregateFromModel(m *models.UserModel) (*aggregates.UserAggregate, error) {
	var middleName string
	if m.MiddleName.Valid {
		middleName = m.MiddleName.String
	}

	var emailVerifiedAt *time.Time
	if m.EmailVerifiedAt.Valid {
		emailVerifiedAt = &m.EmailVerifiedAt.Time
	}

	var passwordChangedAt *time.Time
	if m.PasswordChangedAt.Valid {
		passwordChangedAt = &m.PasswordChangedAt.Time
	}

	var passwordExpiresAt *time.Time
	if m.PasswordExpiresAt.Valid {
		passwordExpiresAt = &m.PasswordExpiresAt.Time
	}

	var lockedUntil *time.Time
	if m.LockedUntil.Valid {
		lockedUntil = &m.LockedUntil.Time
	}

	var lastLoginAt *time.Time
	if m.LastLoginAt.Valid {
		lastLoginAt = &m.LastLoginAt.Time
	}

	var lastActiveAt *time.Time
	if m.LastActiveAt.Valid {
		lastActiveAt = &m.LastActiveAt.Time
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	var lastLoginIP string
	if m.LastLoginIP.Valid {
		lastLoginIP = m.LastLoginIP.String
	}

	var phone string
	if m.Phone.Valid {
		phone = m.Phone.String
	}

	return aggregates.RehydrateUser(
		m.ID,
		m.Email,
		m.PasswordHash,
		m.Status,
		m.FirstName,
		m.LastName,
		middleName,
		lastLoginIP,
		phone,
		m.EmailVerified,
		emailVerifiedAt,
		passwordChangedAt,
		passwordExpiresAt,
		lockedUntil,
		lastLoginAt,
		lastActiveAt,
		deletedAt,
		m.FailedLoginAttempts,
		m.CreatedAt,
		m.UpdatedAt,
		0,
	)
}