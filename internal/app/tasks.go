package app

import "context"

const (
	ErrTaskAlreadyStarted = Error("task already started")
	ErrTaskNotStarted     = Error("task not started")
)

type TaskStarter interface {
	Start(ctx context.Context, taskName string) error
}

type TaskFinisher interface {
	Finish(ctx context.Context, taskName string) error
}
