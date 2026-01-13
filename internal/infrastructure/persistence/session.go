package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	events "github.com/victorotene80/authentication_api/internal/domain/events/types"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Save(ctx context.Context, s *aggregates.SessionAggregate) error {
	tx, err := GetTx(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE sessions
		SET token_hash = $1,
			previous_token_hash = $2,
			rotation_id = $3,
			role = $4,
			status = $5,
			last_seen_at = $6,
			revoked_at = $7,
			expires_at = $8,
			version = version + 1
		WHERE id = $9 AND version = $10
	`

	prevHash := sql.NullString{}
	if s.PreviousTokenHash() != nil {
		prevHash = sql.NullString{String: s.PreviousTokenHash().Value(), Valid: true}
	}

	rotationID := sql.NullString{}
	if s.RotationID() != nil {
		rotationID = sql.NullString{String: *s.RotationID(), Valid: true}
	}

	res, err := tx.ExecContext(ctx, query,
		s.TokenHash().Value(),
		prevHash,
		rotationID,
		s.Role().String(), // role
		string(s.Status()),
		s.LastSeenAt(),
		s.RevokedAt(),
		s.ExpiresAt(),
		s.ID(),
		s.Version(),
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return infrastructure.ErrConflict
	}

	s.CommitVersion()
	return nil
}

func (r *PostgresSessionRepository) FindByID(ctx context.Context, id string) (*aggregates.SessionAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, role, token_hash, previous_token_hash, rotation_id,
		       ip_address, device_id, user_agent,
		       status, created_at, last_seen_at, expires_at, revoked_at, version
		FROM sessions
		WHERE id = $1
	`

	var m models.Session
	err = tx.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.UserID, &m.Role, &m.TokenHash, &m.PreviousTokenHash, &m.RotationID,
		&m.IPAddress, &m.DeviceID, &m.UserAgent,
		&m.Status, &m.CreatedAt, &m.LastSeenAt, &m.ExpiresAt, &m.RevokedAt, &m.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

/*func (r *PostgresSessionRepository) RotateSessionToken(
	ctx context.Context,
	sessionID string,
	newToken valueobjects.SessionTokenHash,
	rotationID string,
	now time.Time,
) (*aggregates.SessionAggregate, error) {
	session, err := r.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := session.RotateKey(newToken, rotationID, now); err != nil {
		return nil, err
	}

	if err := r.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}*/

func (r *PostgresSessionRepository) RotateSessionToken(ctx context.Context, sessionID string, newToken valueobjects.SessionTokenHash, rotationID string, now time.Time) (*aggregates.SessionAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	session, err := r.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	oldToken := session.TokenHash()

	if err := session.RotateKey(newToken, rotationID, now); err != nil {
		return nil, err
	}

	prevHash := sql.NullString{
		String: oldToken.Value(),
		Valid:  true,
	}

	rotation := sql.NullString{
		String: rotationID,
		Valid:  true,
	}

	res, err := tx.ExecContext(ctx, `
		UPDATE sessions
		SET token_hash = $1,
		    previous_token_hash = $2,
		    rotation_id = $3,
		    last_seen_at = $4,
		    version = version + 1
		WHERE id = $5 AND version = $6
	`,
		newToken.Value(),
		prevHash,
		rotation,
		now,
		session.ID(),
		session.Version(),
	)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, infrastructure.ErrConflict
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO session_rotations
		    (session_id, old_token_hash, new_token_hash, rotation_id, rotated_at)
		VALUES ($1, $2, $3, $4, $5)
	`,
		session.ID(),
		oldToken.Value(),
		newToken.Value(),
		rotationID,
		now,
	)
	if err != nil {
		return nil, err
	}

	session.CommitVersion()
	return session, nil
}


func (r *PostgresSessionRepository) RevokeByID(ctx context.Context, sessionID string, now time.Time) (*aggregates.SessionAggregate, error) {
	session, err := r.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := session.Revoke(now, "manual revoke"); err != nil {
		return nil, err
	}

	if err := r.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *PostgresSessionRepository) RevokeAllForUser(ctx context.Context, userID string, now time.Time) ([]*aggregates.SessionAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE sessions
		SET status = 'REVOKED',
		    revoked_at = $1,
		    version = version + 1
		WHERE user_id = $2 AND status = 'ACTIVE'
		RETURNING id, user_id, role, token_hash, previous_token_hash, rotation_id,
		          ip_address, device_id, user_agent,
		          status, created_at, last_seen_at, expires_at, revoked_at, version
	`

	rows, err := tx.QueryContext(ctx, query, now, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revokedSessions []*aggregates.SessionAggregate
	for rows.Next() {
		var s models.Session
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.Role, &s.TokenHash, &s.PreviousTokenHash, &s.RotationID,
			&s.IPAddress, &s.DeviceID, &s.UserAgent,
			&s.Status, &s.CreatedAt, &s.LastSeenAt, &s.ExpiresAt, &s.RevokedAt, &s.Version,
		); err != nil {
			return nil, err
		}

		agg, err := mapSessionModelToAggregate(s)
		if err != nil {
			return nil, err
		}

		agg.RaiseEvent(events.NewSessionRevokedEvent(agg.ID(), agg.UserID(), "batch revoke"))
		revokedSessions = append(revokedSessions, agg)
	}

	return revokedSessions, nil
}

func (r *PostgresSessionRepository) FindActiveByKeyHash(
	ctx context.Context,
	tokenHash valueobjects.SessionTokenHash,
	now time.Time,
) (*aggregates.SessionAggregate, error) {

	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, token_hash, previous_token_hash, rotation_id,
		       ip_address, device_id, user_agent,
		       status, created_at, last_seen_at, expires_at, revoked_at, version
		FROM sessions
		WHERE token_hash = $1
		  AND status = 'ACTIVE'
		  AND expires_at > $2
	`

	var m models.Session
	err = tx.QueryRowContext(ctx, query, tokenHash.Value(), now).Scan(
		&m.ID,
		&m.UserID,
		&m.TokenHash,
		&m.PreviousTokenHash,
		&m.RotationID,
		&m.IPAddress,
		&m.DeviceID,
		&m.UserAgent,
		&m.Status,
		&m.CreatedAt,
		&m.LastSeenAt,
		&m.ExpiresAt,
		&m.RevokedAt,
		&m.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

func (r *PostgresSessionRepository) FindByTokenHash(ctx context.Context, tokenHash valueobjects.SessionTokenHash, now time.Time) (*aggregates.SessionAggregate, error) {
	tx, err := GetTx(ctx)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, user_id, role, token_hash, previous_token_hash, rotation_id,
		       ip_address, device_id, user_agent,
		       status, created_at, last_seen_at, expires_at, revoked_at, version
		FROM sessions
		WHERE (token_hash = $1 OR previous_token_hash = $1)
		  AND expires_at > $2
	`

	var m models.Session
	err = tx.QueryRowContext(ctx, query, tokenHash.Value(), now).Scan(
		&m.ID, &m.UserID, &m.Role, &m.TokenHash, &m.PreviousTokenHash, &m.RotationID,
		&m.IPAddress, &m.DeviceID, &m.UserAgent,
		&m.Status, &m.CreatedAt, &m.LastSeenAt, &m.ExpiresAt, &m.RevokedAt, &m.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infrastructure.ErrSessionNotFound
		}
		return nil, err
	}

	return mapSessionModelToAggregate(m)
}

func mapSessionModelToAggregate(m models.Session) (*aggregates.SessionAggregate, error) {
	roleVO, err := valueobjects.NewRole(m.Role)
	if err != nil {
		return nil, err
	}

	var prevHashVO *valueobjects.SessionTokenHash
	if m.PreviousTokenHash != nil {
		h, err := valueobjects.NewSessionTokenHash(*m.PreviousTokenHash)
		if err != nil {
			return nil, err
		}
		prevHashVO = &h
	}

	return aggregates.RehydrateSession(
		m.ID,
		m.UserID,
		roleVO, 
		m.TokenHash,
		prevHashVOString(prevHashVO),
		m.RotationID,
		m.IPAddress,
		m.UserAgent,
		m.DeviceID,
		m.Status,
		m.CreatedAt,
		m.LastSeenAt,
		m.ExpiresAt,
		m.RevokedAt,
		m.Version,
	)
}

func prevHashVOString(vo *valueobjects.SessionTokenHash) *string {
	if vo == nil {
		return nil
	}
	val := vo.Value()
	return &val
}
