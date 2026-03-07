package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	appmsg "github.com/victorotene80/authentication_api/internal/application/messaging"
	status "github.com/victorotene80/authentication_api/internal/application/outbox"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

type PostgresOutboxRepository struct {
	db *sql.DB
}

func NewPostgresOutboxRepository(db *sql.DB) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{db: db}
}

func (r *PostgresOutboxRepository) Add(ctx context.Context, envelope appmsg.Envelope) error {
	exec := ChooseExecutor(ctx, r.db)

	metadataJSON, err := json.Marshal(envelope.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	model := models.OutboxModel{
		ID:            envelope.ID,
		AggregateID:   envelope.AggregateID,
		EventName:     envelope.Name,
		Payload:       envelope.Payload,
		OccurredAt:    envelope.OccurredAt,
		Version:       envelope.Version,
		Metadata:      metadataJSON,
		AggregateType: envelope.AggregateType,
		Status:        string(status.Pending),
	}

	query := `
		INSERT INTO auth.outbox_events
		(id, aggregate_id, event_name, payload, occurred_at, version, metadata, aggregate_type, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err = exec.ExecContext(ctx, query,
		model.ID,
		model.AggregateID,
		model.EventName,
		model.Payload,
		model.OccurredAt,
		model.Version,
		model.Metadata,
		model.AggregateType,
		model.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert outbox event: %w", err)
	}

	return nil
}

func (r *PostgresOutboxRepository) FetchUnprocessed(ctx context.Context, limit int) ([]appmsg.Envelope, error) {
	query := `
		SELECT id, aggregate_id, event_name, payload, occurred_at, version, metadata, aggregate_type
		FROM auth.outbox_events
		WHERE status = 'PENDING'
		ORDER BY occurred_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envelopes []appmsg.Envelope
	for rows.Next() {
		var model models.OutboxModel
		if err := rows.Scan(
			&model.ID,
			&model.AggregateID,
			&model.EventName,
			&model.Payload,
			&model.OccurredAt,
			&model.Version,
			&model.Metadata,
			&model.AggregateType,
		); err != nil {
			return nil, err
		}

		var meta map[string]string
		if len(model.Metadata) > 0 {
			if err := json.Unmarshal(model.Metadata, &meta); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		} else {
			meta = make(map[string]string)
		}

		kind := appmsg.Kind(meta["message_kind"])
		if kind == "" {
			kind = appmsg.KindIntegrationEvent
		}

		envelopes = append(envelopes, appmsg.Envelope{
			ID:            model.ID,
			Name:          model.EventName,
			Kind:          kind,
			AggregateID:   model.AggregateID,
			AggregateType: model.AggregateType,
			OccurredAt:    model.OccurredAt,
			Payload:       model.Payload,
			Metadata:      meta,
			CorrelationID: meta["correlation_id"],
			CausationID:   meta["causation_id"],
			Version:       model.Version,
		})
	}

	return envelopes, nil
}

func (r *PostgresOutboxRepository) MarkInProgress(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE auth.outbox_events SET status = $1, updated_at = NOW() WHERE id = $2`,
		status.InProgress,
		id,
	)
	return err
}

func (r *PostgresOutboxRepository) MarkSent(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE auth.outbox_events SET status = $1, processed_at = $2, updated_at = NOW() WHERE id = $3`,
		status.Sent,
		time.Now().UTC(),
		id,
	)
	return err
}

func (r *PostgresOutboxRepository) MarkFailed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE auth.outbox_events SET status = $1, updated_at = NOW() WHERE id = $2`,
		status.Failed,
		id,
	)
	return err
}

func (r *PostgresOutboxRepository) ReclaimStaleInProgress(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int, error) {
	const q = `
		WITH stale AS (
			SELECT id
			FROM auth.outbox_events
			WHERE status = 'IN_PROGRESS'
			  AND updated_at < $1
			ORDER BY updated_at
			LIMIT $2
			FOR UPDATE SKIP LOCKED
		)
		UPDATE auth.outbox_events o
		SET status = 'PENDING',
		    updated_at = NOW()
		FROM stale
		WHERE o.id = stale.id
	`
	res, err := r.db.ExecContext(ctx, q, olderThan, limit)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}