package tasks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type EventStore interface {
	Store(ctx context.Context, event Event) error
}

type Starter struct {
	eventStore EventStore
	now        func() time.Time
}

func NewStarter(eventStore EventStore) Starter {
	return Starter{eventStore: eventStore, now: time.Now}
}

func (s Starter) Start(ctx context.Context, taskName string) error {
	event := Event{
		ID:        uuid.New(),
		Type:      EventTypeTaskStarted,
		CreatedAt: s.now(),
		Payload: EventPayload{
			TaskName: taskName,
		},
	}
	if err := s.eventStore.Store(ctx, event); err != nil {
		return fmt.Errorf("storing event: %w", err)
	}
	return nil
}
