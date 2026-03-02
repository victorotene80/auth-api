package persistence

import (
	"context"
	"database/sql"
	"encoding/json"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
)

type PostgresAuditLogger struct {
	db *sql.DB
}

var _ appContracts.AuditLogger = (*PostgresAuditLogger)(nil)

func NewPostgresAuditLogger(db *sql.DB) *PostgresAuditLogger {
	return &PostgresAuditLogger{db: db}
}

func (l *PostgresAuditLogger) Log(
	ctx context.Context,
	rec dto.AuditRecord,
) error {
	const q = `
		INSERT INTO auth.audit_log (
			action,
			user_id,
			actor_id,
			api_key_id,
			session_id,
			organization_id,
			ip_address,
			user_agent,
			country_code,
			target_resource,
			target_id,
			metadata,
			success,
			failure_reason
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12,
			$13, $14
		)
	`

	var metadataJSON any
	if rec.Metadata != nil {
		b, err := json.Marshal(rec.Metadata)
		if err != nil {
			metadataJSON = nil
		} else {
			metadataJSON = b
		}
	} else {
		metadataJSON = nil
	}

	_, err := l.db.ExecContext(
		ctx,
		q,
		rec.Action,
		nullableString(rec.UserID),
		nullableString(rec.ActorID),
		nullableString(rec.APIKeyID),
		nullableString(rec.SessionID),
		nullableString(rec.OrganizationID),
		nullableString(rec.IPAddress),
		nullableString(rec.UserAgent),
		nullableString(rec.CountryCode),
		nullableString(rec.TargetResource),
		nullableString(rec.TargetID),
		metadataJSON,
		rec.Success,
		nullableString(rec.FailureReason),
	)

	return err
}

func nullableString(s *string) any {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}
