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

// Event represents something that has happened in relation to a task.
type Event struct {
	ID        uuid.UUID
	Type      EventType
	TaskName  string
	CreatedAt time.Time
}

//go:generate mockery --name=EventStore
// EventStore is used to store individual events related to tasks.
type EventStore interface {
	Store(ctx context.Context, event Event) error
}

//go:generate mockery --name=EventFinder
// EventFinder is used to find specific events related to tasks. It will return ErrEventNotFound if there is no
// event with the given name.
type EventFinder interface {
	LatestByName(ctx context.Context, taskName string) (Event, error)
}
