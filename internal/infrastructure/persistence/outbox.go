package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/victorotene80/authentication_api/internal/application/messaging"
	status "github.com/victorotene80/authentication_api/internal/application/outbox"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/models"
)

type PostgresOutboxRepository struct {
	db *sql.DB
}

func NewPostgresOutboxRepository(db *sql.DB) *PostgresOutboxRepository {
	return &PostgresOutboxRepository{db: db}
}

func (r *PostgresOutboxRepository) Add(ctx context.Context, envelope messaging.Envelope) error {
	metadataJSON, err := json.Marshal(envelope.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	model := models.OutboxModel{
		ID:            envelope.ID,
		AggregateID:   envelope.AggregateID,
		EventName:     envelope.EventName,
		Payload:       envelope.Payload,
		OccurredAt:    envelope.OccurredAt,
		Version:       envelope.Version,
		Metadata:      metadataJSON,
		AggregateType: envelope.Metadata["aggregate_type"],
		Status:        string(status.Pending),
	}

	query := `
		INSERT INTO outbox_events
		(id, aggregate_id, event_name, payload, occurred_at, version, metadata, aggregate_type, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err = r.db.ExecContext(ctx, query,
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

func (r *PostgresOutboxRepository) FetchUnprocessed(ctx context.Context, limit int) ([]messaging.Envelope, error) {
	query := `
		SELECT id, aggregate_id, event_name, payload, occurred_at, version, metadata, aggregate_type
		FROM outbox_events
		WHERE status = 'PENDING'
		ORDER BY occurred_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envelopes []messaging.Envelope
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
		meta["aggregate_type"] = model.AggregateType

		envelopes = append(envelopes, messaging.Envelope{
			ID:          model.ID,
			AggregateID: model.AggregateID,
			EventName:   model.EventName,
			Payload:     model.Payload,
			OccurredAt:  model.OccurredAt,
			Version:     model.Version,
			Metadata:    meta,
		})
	}

	return envelopes, nil
}

func (r *PostgresOutboxRepository) MarkInProgress(ctx context.Context, id string) error {
	query := `UPDATE outbox_events SET status=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, status.InProgress, id)
	return err
}

func (r *PostgresOutboxRepository) MarkSent(ctx context.Context, id string) error {
	query := `UPDATE outbox_events SET status=$1, processed_at=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, query, status.Sent, time.Now().UTC(), id)
	return err
}

func (r *PostgresOutboxRepository) MarkFailed(ctx context.Context, id string) error {
	query := `UPDATE outbox_events SET status=$1 WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, status.Failed, id)
	return err
}
