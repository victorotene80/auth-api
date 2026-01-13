package models

import "time"

type OutboxModel struct {
	ID            string     `db:"id"`
	AggregateID   string     `db:"aggregate_id"`
	EventName     string     `db:"event_name"`
	Payload       []byte     `db:"payload"`
	OccurredAt    time.Time  `db:"occurred_at"`
	Version       int        `db:"version"`
	Metadata      []byte     `db:"metadata"`      
	AggregateType string     `db:"aggregate_type"`
	ProcessedAt   *time.Time `db:"processed_at"`
	Status        string     `db:"status"`
}
