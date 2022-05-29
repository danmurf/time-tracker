package tasks

import (
	"context"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
	"github.com/google/uuid"
	"time"
)

type Starter struct {
	eventStore app.EventStore
	now        func() time.Time
}

func NewStarter(eventStore app.EventStore) Starter {
	return Starter{eventStore: eventStore, now: time.Now}
}

func (s Starter) Start(ctx context.Context, taskName string) error {
	event := app.Event{
		ID:        uuid.New(),
		Type:      app.EventTypeTaskStarted,
		TaskName:  taskName,
		CreatedAt: s.now(),
	}
	if err := s.eventStore.Store(ctx, event); err != nil {
		return fmt.Errorf("storing event: %w", err)
	}
	return nil
}
