package app

import (
	"context"
	"github.com/google/uuid"
	"time"
)

const (
	EventTypeTaskStarted  = EventType("task-started")
	EventTypeTaskFinished = EventType("task-finished")

	ErrEventNotFound = Error("event not found")
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

//go:generate mockery --name=EventFinder
type EventFinder interface {
	LatestByName(ctx context.Context, taskName string) (Event, error)
}
