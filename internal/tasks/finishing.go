package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
	"github.com/google/uuid"
	"time"
)

var _ app.TaskFinisher = (*Finisher)(nil)

type Finisher struct {
	eventStore  app.EventStore
	eventFinder app.EventFinder
	now         func() time.Time
	newUUID     func() uuid.UUID
}

func NewFinisher(eventStore app.EventStore, eventFinder app.EventFinder) Finisher {
	return Finisher{eventStore: eventStore, eventFinder: eventFinder, now: time.Now, newUUID: uuid.New}
}

func (f Finisher) Finish(ctx context.Context, taskName string) error {
	latest, err := f.eventFinder.LatestByName(ctx, taskName)
	switch {
	case err != nil && !errors.Is(err, app.ErrEventNotFound):
		return fmt.Errorf("finding latest event: %w", err)
	case !errors.Is(err, app.ErrEventNotFound) && latest.Type == app.EventTypeTaskFinished:
		return fmt.Errorf("task previously finished: %w", app.ErrTaskNotStarted)
	case errors.Is(err, app.ErrEventNotFound):
		return fmt.Errorf("task never started: %w", app.ErrTaskNotStarted)
	}

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
