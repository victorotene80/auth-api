package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	events "github.com/victorotene80/authentication_api/internal/domain/events/types"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

var _ repository.SessionRepository = (*PostgresSessionRepository)(nil)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func hashPtrToStringPtr(h *valueobjects.SessionTokenHash) *string {
	if h == nil {
		return nil
	}
	v := h.Value()
	return &v
}

func (r *PostgresSessionRepository) Save(
	ctx context.Context,
	s *aggregates.SessionAggregate,
) error {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
INSERT INTO auth.sessions (
    id,
    user_id,
    token_hash,
    refresh_token_hash,
    ip_address,
    user_agent,
    device_fingerprint,
    device_name,
    country_code,
    city,
    is_mfa_verified,
    impersonated_by,
    last_active_at,
    expires_at,
    revoked_at,
    revoke_reason,
    created_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8,
    $9, $10, $11, $12,
    $13, $14, $15, $16,
    $17
)
ON CONFLICT (id) DO UPDATE SET
    token_hash        = EXCLUDED.token_hash,
    refresh_token_hash = EXCLUDED.refresh_token_hash,
    ip_address        = EXCLUDED.ip_address,
    user_agent        = EXCLUDED.user_agent,
    device_fingerprint = EXCLUDED.device_fingerprint,
    device_name       = EXCLUDED.device_name,
    country_code      = EXCLUDED.country_code,
    city              = EXCLUDED.city,
    is_mfa_verified   = EXCLUDED.is_mfa_verified,
    impersonated_by   = EXCLUDED.impersonated_by,
    last_active_at    = EXCLUDED.last_active_at,
    expires_at        = EXCLUDED.expires_at,
    revoked_at        = EXCLUDED.revoked_at,
    revoke_reason     = EXCLUDED.revoke_reason
`

	_, err := exec.ExecContext(ctx, q,
		s.ID(),
		s.UserID(),
		s.TokenHash().Value(),
		hashPtrToStringPtr(s.RefreshTokenHash()),
		stringPtrOrNil(s.IPAddress()),
		stringPtrOrNil(s.UserAgent()),
		stringPtrOrNil(s.DeviceFingerprint()),
		stringPtrOrNil(s.DeviceName()),
		stringPtrOrNil(s.CountryCode()),
		stringPtrOrNil(s.City()),
		s.IsMFAVerified(),
		s.ImpersonatedBy(),
		s.LastActiveAt(),
		s.ExpiresAt(),
		s.RevokedAt(),
		s.RevokeReason(),
		s.CreatedAt(),
	)
	return err
}

func (r *PostgresSessionRepository) FindByID(
	ctx context.Context,
	id string,
) (*aggregates.SessionAggregate, error) {
	exec := ChooseExecutor(ctx, r.db)

	const q = `
SELECT
    id,
    user_id,
    token_hash,
    refresh_token_hash,
    ip_address,
    user_agent,
    device_fingerprint,
    device_name,
    country_code,
    city,
    is_mfa_verified,
    impersonated_by,
    last_active_at,
    expires_at,
    revoked_at,
    revoke_reason,
    created_at
FROM auth.sessions
WHERE id = $1
`

	var m models.Session
	if err := exec.QueryRowContext(ctx, q, id).Scan(
		&m.ID,
		&m.UserID,
		&m.TokenHash,
		&m.RefreshTokenHash,
		&m.IPAddress,
		&m.UserAgent,
		&m.DeviceFingerprint,
		&m.DeviceName,
		&m.CountryCode,
		&m.City,
		&m.IsMFAVerified,
		&m.ImpersonatedBy,
		&m.LastActiveAt,
		&m.ExpiresAt,
		&m.RevokedAt,
		&m.RevokeReason,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

func (r *PostgresSessionRepository) RotateSessionToken(
	ctx context.Context,
	sessionID string,
	newToken valueobjects.SessionTokenHash,
	_ string, // rotationID unused now (DB has no place to store it)
	now time.Time,
) (*aggregates.SessionAggregate, error) {
	exec := ChooseExecutor(ctx, r.db)

	session, err := r.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if !session.IsValid(now) {
		return nil, infrastructure.ErrSessionNotFound
	}

	session.RotateKey(newToken, now)

	const q = `
		UPDATE auth.sessions
		SET token_hash = $1,
			last_active_at = $2
		WHERE id = $3
		`

	if _, err := exec.ExecContext(ctx, q,
		newToken.Value(),
		now,
		sessionID,
	); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *PostgresSessionRepository) RevokeByID(
	ctx context.Context,
	sessionID string,
	now time.Time,
) (*aggregates.SessionAggregate, error) {

	exec := ChooseExecutor(ctx, r.db)

	session, err := r.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	session.Revoke(now, "manual revoke")

	const q = `
UPDATE auth.sessions
SET revoked_at = $1,
    revoke_reason = $2
WHERE id = $3
`

	if _, err := exec.ExecContext(ctx, q,
		session.RevokedAt(),
		session.RevokeReason(),
		sessionID,
	); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *PostgresSessionRepository) RevokeAllForUser(
	ctx context.Context,
	userID string,
	now time.Time,
) ([]*aggregates.SessionAggregate, error) {

	exec := ChooseExecutor(ctx, r.db)

	const reason = "batch revoke"

	const q = `
UPDATE auth.sessions
SET revoked_at = $1,
    revoke_reason = $2
WHERE user_id = $3
  AND revoked_at IS NULL
  AND expires_at > $1
RETURNING
    id,
    user_id,
    token_hash,
    refresh_token_hash,
    ip_address,
    user_agent,
    device_fingerprint,
    device_name,
    country_code,
    city,
    is_mfa_verified,
    impersonated_by,
    last_active_at,
    expires_at,
    revoked_at,
    revoke_reason,
    created_at
`

	rows, err := exec.QueryContext(ctx, q, now, reason, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revokedSessions []*aggregates.SessionAggregate

	for rows.Next() {
		var m models.Session
		if err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.TokenHash,
			&m.RefreshTokenHash,
			&m.IPAddress,
			&m.UserAgent,
			&m.DeviceFingerprint,
			&m.DeviceName,
			&m.CountryCode,
			&m.City,
			&m.IsMFAVerified,
			&m.ImpersonatedBy,
			&m.LastActiveAt,
			&m.ExpiresAt,
			&m.RevokedAt,
			&m.RevokeReason,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}

		agg, err := mapSessionModelToAggregate(m)
		if err != nil {
			return nil, err
		}

		agg.RaiseEvent(events.NewSessionRevokedEvent(agg.ID(), agg.UserID(), reason))
		revokedSessions = append(revokedSessions, agg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return revokedSessions, nil
}

func (r *PostgresSessionRepository) FindActiveByKeyHash(
	ctx context.Context,
	tokenHash valueobjects.SessionTokenHash,
	now time.Time,
) (*aggregates.SessionAggregate, error) {

	exec := ChooseExecutor(ctx, r.db)

	const q = `
SELECT
    id,
    user_id,
    token_hash,
    refresh_token_hash,
    ip_address,
    user_agent,
    device_fingerprint,
    device_name,
    country_code,
    city,
    is_mfa_verified,
    impersonated_by,
    last_active_at,
    expires_at,
    revoked_at,
    revoke_reason,
    created_at
FROM auth.sessions
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > $2
`

	var m models.Session
	if err := exec.QueryRowContext(ctx, q, tokenHash.Value(), now).Scan(
		&m.ID,
		&m.UserID,
		&m.TokenHash,
		&m.RefreshTokenHash,
		&m.IPAddress,
		&m.UserAgent,
		&m.DeviceFingerprint,
		&m.DeviceName,
		&m.CountryCode,
		&m.City,
		&m.IsMFAVerified,
		&m.ImpersonatedBy,
		&m.LastActiveAt,
		&m.ExpiresAt,
		&m.RevokedAt,
		&m.RevokeReason,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

func (r *PostgresSessionRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash valueobjects.SessionTokenHash,
	now time.Time,
) (*aggregates.SessionAggregate, error) {

	exec := ChooseExecutor(ctx, r.db)

	const q = `
SELECT
    id,
    user_id,
    token_hash,
    refresh_token_hash,
    ip_address,
    user_agent,
    device_fingerprint,
    device_name,
    country_code,
    city,
    is_mfa_verified,
    impersonated_by,
    last_active_at,
    expires_at,
    revoked_at,
    revoke_reason,
    created_at
FROM auth.sessions
WHERE token_hash = $1
  AND expires_at > $2
`

	var m models.Session
	if err := exec.QueryRowContext(ctx, q, tokenHash.Value(), now).Scan(
		&m.ID,
		&m.UserID,
		&m.TokenHash,
		&m.RefreshTokenHash,
		&m.IPAddress,
		&m.UserAgent,
		&m.DeviceFingerprint,
		&m.DeviceName,
		&m.CountryCode,
		&m.City,
		&m.IsMFAVerified,
		&m.ImpersonatedBy,
		&m.LastActiveAt,
		&m.ExpiresAt,
		&m.RevokedAt,
		&m.RevokeReason,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

func mapSessionModelToAggregate(m models.Session) (*aggregates.SessionAggregate, error) {
	var refreshHash *string
	if m.RefreshTokenHash != nil {
		refreshHash = m.RefreshTokenHash
	}

	ip := ""
	if m.IPAddress != nil {
		ip = *m.IPAddress
	}

	ua := ""
	if m.UserAgent != nil {
		ua = *m.UserAgent
	}

	fp := ""
	if m.DeviceFingerprint != nil {
		fp = *m.DeviceFingerprint
	}

	deviceName := ""
	if m.DeviceName != nil {
		deviceName = *m.DeviceName
	}

	country := ""
	if m.CountryCode != nil {
		country = *m.CountryCode
	}

	city := ""
	if m.City != nil {
		city = *m.City
	}

	var revokedAt *time.Time
	if m.RevokedAt != nil {
		revokedAt = m.RevokedAt
	}

	var revokeReason *string
	if m.RevokeReason != nil {
		revokeReason = m.RevokeReason
	}

	return aggregates.RehydrateSession(
		m.ID,
		m.UserID,
		m.TokenHash,
		refreshHash,
		ip,
		ua,
		fp,
		deviceName,
		country,
		city,
		m.IsMFAVerified,
		m.ImpersonatedBy,
		m.CreatedAt,
		m.LastActiveAt,
		m.ExpiresAt,
		revokedAt,
		revokeReason,
		0,
	)
}
