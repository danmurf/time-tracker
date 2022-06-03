package tasks

import (
	"context"
	"github.com/danmurf/time-tracker/internal/app"
)

var _ app.LastCompletedFetcher = (*Durations)(nil)

type Durations struct {
}

func (d Durations) FetchLastCompleted(ctx context.Context, taskName string) (app.CompletedTask, error) {
	//TODO implement me
	panic("implement me")
}
