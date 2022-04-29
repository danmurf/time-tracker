package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	eventStoreCreation = `
CREATE TABLE IF NOT EXISTS "event_store" (
	"id" varchar NOT NULL,
	"type" varchar NOT NULL DEFAULT NULL,
	"created_at" datetime NOT NULL,
	"payload" text DEFAULT NULL, 
	PRIMARY KEY (id)
);
`
	EventTypeTaskStarted  = "task-started"
	EventTypeTaskFinished = "task-finished"
)

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
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return fmt.Errorf("encoding payload: %w", err)
	}

	if _, err = s.db.ExecContext(ctx, "INSERT INTO `event_store` VALUES(?, ?, ?, ?);", e.ID, e.Type, e.CreatedAt, payload); err != nil {
		return fmt.Errorf("inserting into db: %w", err)
	}

	return nil
}

func (s SQLEventStore) FetchAll(ctx context.Context) ([]Event, error) {
	var events []Event
	rows, err := s.db.QueryContext(ctx, "SELECT id, type, created_at, payload FROM `event_store` ORDER BY created_at DESC;")
	if err != nil {
		return events, fmt.Errorf("querying db: %w", err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		return events, fmt.Errorf("reading rows: %w", err)
	}
	for rows.Next() {
		var event Event
		var payload []byte
		if err = rows.Scan(&event.ID, &event.Type, &event.CreatedAt, &payload); err != nil {
			return events, fmt.Errorf("scanning row: %w", err)
		}
		if err := json.Unmarshal(payload, &event.Payload); err != nil {
			return events, fmt.Errorf("decoding payload: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (s SQLEventStore) bootstrap(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, eventStoreCreation)
	if err != nil {
		return fmt.Errorf("creating event store: %w", err)
	}

	return nil
}
