package events

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const eventStoreCreation = `
CREATE TABLE IF NOT EXISTS "event_store" (
	"id" varchar NOT NULL,
	"type" varchar NOT NULL DEFAULT NULL,
	"created_at" datetime NOT NULL,
	"payload" text DEFAULT NULL, 
	PRIMARY KEY (id)
);
`

type EventType string

type Event struct {
	ID        uuid.UUID
	Type      EventType
	CreatedAt time.Time
	Payload   EventPayload
}

type EventPayload struct {
	TaskName string `json:"task_name"`
}

func NewSQLEventStore(ctx context.Context, db *sql.DB) (SQLEventStore, error) {
	s := SQLEventStore{db: db}
	if err := s.bootstrap(ctx); err != nil {
		return SQLEventStore{}, fmt.Errorf("bootstrapping sql event store: %w", err)
	}
	return s, nil
}

type SQLEventStore struct {
	db *sql.DB
}

func (s SQLEventStore) Store(ctx context.Context, e Event) error {
	panic("implement me")
}

func (s SQLEventStore) Fetch(ctx context.Context, e Event) error {
	panic("implement me")
}

func (s SQLEventStore) bootstrap(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, eventStoreCreation)
	if err != nil {
		return fmt.Errorf("creating event store: %w", err)
	}
	return nil
}
