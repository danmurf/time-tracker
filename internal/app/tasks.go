package app

import (
	"context"
	"time"
)

const (
	ErrTaskAlreadyStarted = Error("task already started")
	ErrTaskNotStarted     = Error("task not started")
	ErrTaskNeverCompleted = Error("task never completed")
)

// CompletedTask represents a task which has been both started and finished. The CompletedTask.Duration field
// represents the duration difference between CompletedTask.Started and CompletedTask.Finished.
type CompletedTask struct {
	Name     string
	Started  time.Time
	Finished time.Time
	Duration time.Duration
}

// TaskStarter is used to start a task with the given name. It can return ErrTaskAlreadyStarted if the task has
// already been started.
type TaskStarter interface {
	Start(ctx context.Context, taskName string) error
}

// TaskFinisher is used to finish a currently running task. It can return ErrTaskNotStarted if the task is not
// currently in progress.
type TaskFinisher interface {
	Finish(ctx context.Context, taskName string) error
}

// LastCompletedFetcher is used to fetch the last completed task with the given name. It can return
// ErrTaskNeverCompleted if the task has never been completed (started and finished). If a task has been completed,
// and another instance is in progress, the last completed version will be returned.
type LastCompletedFetcher interface {
	FetchLastCompleted(ctx context.Context, taskName string) (CompletedTask, error)
}
