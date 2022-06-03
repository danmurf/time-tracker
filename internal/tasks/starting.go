package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
	"github.com/google/uuid"
	"time"
)

var _ app.TaskStarter = (*Starter)(nil)

type Starter struct {
	eventStore  app.EventStore
	eventFinder app.EventFinder
	now         func() time.Time
	newUUID     func() uuid.UUID
}

func NewStarter(eventStore app.EventStore, eventFinder app.EventFinder) Starter {
	return Starter{eventStore: eventStore, eventFinder: eventFinder, now: time.Now, newUUID: uuid.New}
}

func (s Starter) Start(ctx context.Context, taskName string) error {
	latest, err := s.eventFinder.LatestByName(ctx, taskName)
	switch {
	case err != nil && !errors.Is(err, app.ErrEventNotFound):
		return fmt.Errorf("finding latest event: %w", err)
	case !errors.Is(err, app.ErrEventNotFound) && latest.Type == app.EventTypeTaskStarted:
		return fmt.Errorf("task started event found: %w", app.ErrTaskAlreadyStarted)
	}

	if err := s.eventStore.Store(ctx, app.Event{
		ID:        s.newUUID(),
		Type:      app.EventTypeTaskStarted,
		TaskName:  taskName,
		CreatedAt: s.now(),
	}); err != nil {
		return fmt.Errorf("storing event: %w", err)
	}
	return nil
}
