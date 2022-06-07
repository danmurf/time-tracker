package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
)

var _ app.LastCompletedFetcher = (*Durations)(nil)

type Durations struct {
	eventFinder app.EventFinder
}

func NewDurations(eventFinder app.EventFinder) Durations {
	return Durations{eventFinder: eventFinder}
}

func (d Durations) FetchLastCompleted(ctx context.Context, taskName string) (ct app.CompletedTask, err error) {
	finished, err := d.eventFinder.LatestByNameType(ctx, taskName, app.EventTypeTaskFinished)
	switch {
	case err != nil && !errors.Is(err, app.ErrEventNotFound):
		return ct, fmt.Errorf("finding finished event: %w", err)
	case errors.Is(err, app.ErrEventNotFound):
		return ct, fmt.Errorf("finding finished event: %w", app.ErrTaskNeverCompleted)
	}

	started, err := d.eventFinder.LatestByNameType(ctx, taskName, app.EventTypeTaskStarted)
	switch {
	case err != nil:
		// Started event should always exist, if there is a finished event
		return ct, fmt.Errorf("finding started event: %w", err)
	}

	return app.CompletedTask{
		Name:     taskName,
		Started:  started.CreatedAt,
		Finished: finished.CreatedAt,
		Duration: finished.CreatedAt.Sub(started.CreatedAt),
	}, nil
}
