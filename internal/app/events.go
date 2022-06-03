package app

import (
	"context"
	"github.com/google/uuid"
	"time"
)

const (
	EventTypeTaskStarted  = EventType("task-started")
	EventTypeTaskFinished = EventType("task-finished")
)

type EventType string

type Event struct {
	ID        uuid.UUID
	Type      EventType
	TaskName  string
	CreatedAt time.Time
}

//go:generate mockery --name=EventStore
type EventStore interface {
	Store(ctx context.Context, event Event) error
}
