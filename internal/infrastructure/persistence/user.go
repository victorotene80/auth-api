package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, agg *aggregates.UserAggregate) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (
			id,
			email,
			password_hash,
			role,
			first_name,
			last_name,
			middle_name,
			is_active,
			is_verified,
			last_login_at,
			failed_login_attempts,
			last_failed_at,
			locked_at,
			created_at,
			updated_at,
			version
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16
		)
	`

	_, err = tx.ExecContext(ctx, query,
		agg.User.ID(),
		agg.User.Email().String(),
		agg.User.Password().Value(),
		agg.User.Role().String(),
		agg.User.FirstName,
		agg.User.LastName,
		agg.User.MiddleName,
		agg.User.IsActive(),
		agg.User.IsVerified(),
		agg.User.LastLoginAt(),
		agg.User.FailedLoginAttempts,
		agg.User.LastFailedAt,
		agg.User.LockedAt,
		agg.CreatedAt(),
		agg.UpdatedAt(),
		agg.Version(),
	)
	return err
}

func (r *PostgresUserRepository) Update(ctx context.Context, agg *aggregates.UserAggregate) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE users SET
			email=$1,
			password_hash=$2,
			role=$3,
			first_name=$4,
			last_name=$5,
			middle_name=$6,
			is_active=$7,
			is_verified=$8,
			last_login_at=$9,
			failed_login_attempts=$10,
			last_failed_at=$11,
			locked_at=$12,
			updated_at=$13,
			version=version+1
		WHERE id=$14 AND version=$15

	`

	res, err := tx.ExecContext(ctx, query,
		agg.User.Email().String(),
		agg.User.Password().Value(),
		agg.User.Role().String(),
		agg.User.FirstName,
		agg.User.LastName,
		agg.User.MiddleName,
		agg.User.FailedLoginAttempts,
		agg.User.LastFailedAt,
		agg.User.LockedAt,
		agg.UpdatedAt(),
		agg.User.ID(),
		agg.Version(),
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("optimistic concurrency: user was updated by another process")
	}

	agg.CommitVersion()
	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*aggregates.UserAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id,email,password_hash,role,first_name,last_name,middle_name,failed_login_attempts,last_failed_at,locked_at,created_at,updated_at,version
		FROM users WHERE id=$1
	`

	row := tx.QueryRowContext(ctx, query, id)
	m := &models.UserModel{}
	err = row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.Role,
		&m.FirstName,
		&m.LastName,
		&m.MiddleName,
		&m.FailedLoginAttempts,
		&m.LastFailedAt,
		&m.LockedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return aggregates.RehydrateUser(
		m.ID,
		m.Email,
		m.PasswordHash,
		m.Role,
		m.FirstName,
		m.LastName,
		m.MiddleName,
		m.FailedLoginAttempts,
		m.LastFailedAt,
		m.LockedAt,
		m.CreatedAt,
		m.UpdatedAt,
		m.Version,
	)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*aggregates.UserAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id,email,password_hash,role,first_name,last_name,middle_name,failed_login_attempts,last_failed_at,locked_at,created_at,updated_at,version
		FROM users WHERE email=$1
	`

	row := tx.QueryRowContext(ctx, query, email.String())
	m := &models.UserModel{}
	err = row.Scan(
		&m.ID,
		&m.Email,
		&m.PasswordHash,
		&m.Role,
		&m.FirstName,
		&m.LastName,
		&m.MiddleName,
		&m.FailedLoginAttempts,
		&m.LastFailedAt,
		&m.LockedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return aggregates.RehydrateUser(
		m.ID,
		m.Email,
		m.PasswordHash,
		m.Role,
		m.FirstName,
		m.LastName,
		m.MiddleName,
		m.FailedLoginAttempts,
		m.LastFailedAt,
		m.LockedAt,
		m.CreatedAt,
		m.UpdatedAt,
		m.Version,
	)
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email valueobjects.Email) (bool, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return false, err
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err = tx.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
}

func (r *PostgresUserRepository) List(
	ctx context.Context,
	page, pageSize int,
	role *valueobjects.Role,
	isActive *bool,
) ([]*aggregates.UserAggregate, int64, error) {

	tx, err := GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	// Build query dynamically based on filters
	query := `
		SELECT id,email,password_hash,role,first_name,last_name,middle_name,
		       failed_login_attempts,last_failed_at,locked_at,
		       created_at,updated_at,version
		FROM users
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argIdx := 1

	if role != nil {
		query += fmt.Sprintf(" AND role=$%d", argIdx)
		args = append(args, role.String())
		argIdx++
	}
	if isActive != nil {
		query += fmt.Sprintf(" AND is_active=$%d", argIdx)
		args = append(args, *isActive)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT $%d OFFSET $%d"
	args = append(args, pageSize, offset)
	query = fmt.Sprintf(query, argIdx, argIdx+1)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*aggregates.UserAggregate
	for rows.Next() {
		var m models.UserModel
		if err := rows.Scan(
			&m.ID,
			&m.Email,
			&m.PasswordHash,
			&m.Role,
			&m.FirstName,
			&m.LastName,
			&m.MiddleName,
			&m.FailedLoginAttempts,
			&m.LastFailedAt,
			&m.LockedAt,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Version,
		); err != nil {
			return nil, 0, err
		}

		agg, err := aggregates.RehydrateUser(
			m.ID,
			m.Email,
			m.PasswordHash,
			m.Role,
			m.FirstName,
			m.LastName,
			m.MiddleName,
			m.FailedLoginAttempts,
			m.LastFailedAt,
			m.LockedAt,
			m.CreatedAt,
			m.UpdatedAt,
			m.Version,
		)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, agg)
	}

	// Count total records for pagination
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	var total int64
	if err := tx.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
