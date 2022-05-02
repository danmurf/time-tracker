package tasks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Finisher struct {
	eventStore EventStore
	now        func() time.Time
}

func NewFinisher(eventStore EventStore) Finisher {
	return Finisher{eventStore: eventStore, now: time.Now}
}

func (f Finisher) Finish(ctx context.Context, taskName string) error {
	event := Event{
		ID:        uuid.New(),
		Type:      EventTypeTaskFinished,
		TaskName:  taskName,
		CreatedAt: f.now(),
	}
	if err := f.eventStore.Store(ctx, event); err != nil {
		return fmt.Errorf("storing event: %w", err)
	}
	return nil
}
