package tasks

import (
	"context"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
	"github.com/google/uuid"
	"time"
)

type Finisher struct {
	eventStore app.EventStore
	now        func() time.Time
	newUUID    func() uuid.UUID
}

func NewFinisher(eventStore app.EventStore) Finisher {
	return Finisher{eventStore: eventStore, now: time.Now}
}

func (f Finisher) Finish(ctx context.Context, taskName string) error {
	event := app.Event{
		ID:        f.newUUID(),
		Type:      app.EventTypeTaskFinished,
		TaskName:  taskName,
		CreatedAt: f.now(),
	}
	if err := f.eventStore.Store(ctx, event); err != nil {
		return fmt.Errorf("storing event: %w", err)
	}
	return nil
}
