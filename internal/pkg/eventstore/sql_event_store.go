package eventstore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
)

const (
	eventStoreCreation = `
CREATE TABLE IF NOT EXISTS "event_store" (
	"id" varchar NOT NULL,
	"type" varchar NOT NULL DEFAULT NULL,
	"task_name" varchar NOT NULL DEFAULT NULL, 
	"created_at" datetime NOT NULL,
	PRIMARY KEY (id)
);
`
)

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

func (s SQLEventStore) Store(ctx context.Context, e app.Event) error {
	if _, err := s.db.ExecContext(ctx, "INSERT INTO `event_store` VALUES(?, ?, ?, ?);", e.ID, e.Type, e.TaskName, e.CreatedAt); err != nil {
		return fmt.Errorf("inserting into db: %w", err)
	}

	return nil
}

func (s SQLEventStore) FetchAll(ctx context.Context) ([]app.Event, error) {
	var events []app.Event
	rows, err := s.db.QueryContext(ctx, "SELECT id, type, task_name, created_at FROM `event_store` ORDER BY created_at DESC;")
	if err != nil {
		return events, fmt.Errorf("querying db: %w", err)
	}
	defer rows.Close()

	if rows.Err() != nil {
		return events, fmt.Errorf("reading rows: %w", err)
	}
	for rows.Next() {
		var event app.Event
		if err = rows.Scan(&event.ID, &event.Type, &event.TaskName, &event.CreatedAt); err != nil {
			return events, fmt.Errorf("scanning row: %w", err)
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
