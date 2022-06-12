package tasks

import (
	"context"
	"errors"
	"github.com/danmurf/time-tracker/internal/app"
	app_mocks "github.com/danmurf/time-tracker/internal/app/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestDurations_FetchLastCompleted(t *testing.T) {
	now := time.Now()
	startedID := uuid.New()
	finishedID := uuid.New()
	type fields struct {
		eventFinder *app_mocks.EventFinder
	}
	type args struct {
		ctx      context.Context
		taskName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    app.CompletedTask
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successfully finds completed task (1 minute)",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByNameType", mock.Anything, "test-task", mock.AnythingOfType("app.EventType")).
						Twice().
						Return(func(_ context.Context, taskName string, eventType app.EventType) app.Event {
							createdAt := now.Add(-2 * time.Minute)
							id := startedID
							if eventType == app.EventTypeTaskFinished {
								createdAt = now.Add(-1 * time.Minute)
								id = finishedID
							}
							return app.Event{
								ID:        id,
								Type:      eventType,
								TaskName:  taskName,
								CreatedAt: createdAt,
							}
						}, nil)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test-task",
			},
			want: app.CompletedTask{
				Name:     "test-task",
				Started:  app.Event{ID: startedID, Type: app.EventTypeTaskStarted, TaskName: "test-task", CreatedAt: now.Add(-2 * time.Minute)},
				Finished: app.Event{ID: finishedID, Type: app.EventTypeTaskFinished, TaskName: "test-task", CreatedAt: now.Add(-1 * time.Minute)},
				Duration: 1 * time.Minute,
			},
			wantErr: assert.NoError,
		},
		{
			name: "successfully finds completed task (2 seconds)",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByNameType", mock.Anything, "test-task", mock.AnythingOfType("app.EventType")).
						Twice().
						Return(func(_ context.Context, taskName string, eventType app.EventType) app.Event {
							createdAt := now.Add(-5 * time.Second)
							id := startedID
							if eventType == app.EventTypeTaskFinished {
								createdAt = now.Add(-3 * time.Second)
								id = finishedID
							}
							return app.Event{
								ID:        id,
								Type:      eventType,
								TaskName:  taskName,
								CreatedAt: createdAt,
							}
						}, nil)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test-task",
			},
			want: app.CompletedTask{
				Name:     "test-task",
				Started:  app.Event{ID: startedID, Type: app.EventTypeTaskStarted, TaskName: "test-task", CreatedAt: now.Add(-5 * time.Second)},
				Finished: app.Event{ID: finishedID, Type: app.EventTypeTaskFinished, TaskName: "test-task", CreatedAt: now.Add(-3 * time.Second)},
				Duration: 2 * time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "finds task which has never finished",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByNameType", mock.Anything, "test-task", app.EventTypeTaskFinished).
						Once().
						Return(app.Event{}, app.ErrEventNotFound)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test-task",
			},
			want: app.CompletedTask{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Truef(t, errors.Is(err, app.ErrTaskNeverCompleted), "want err [%s] got [%s]", app.ErrTaskNeverCompleted, err)
			},
		},
		{
			name: "finds task which has finished, but never started (invalid state)",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByNameType", mock.Anything, "test-task", mock.AnythingOfType("app.EventType")).
						Twice().
						Return(func(_ context.Context, taskName string, eventType app.EventType) app.Event {
							if eventType == app.EventTypeTaskStarted {
								return app.Event{}
							}
							return app.Event{
								ID:        uuid.New(),
								Type:      eventType,
								TaskName:  taskName,
								CreatedAt: now.Add(-1 * time.Minute),
							}
						}, func(_ context.Context, taskName string, eventType app.EventType) error {
							if eventType == app.EventTypeTaskFinished {
								return nil
							}
							return app.ErrEventNotFound
						})
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test-task",
			},
			want:    app.CompletedTask{},
			wantErr: assert.Error,
		},
		{
			name: "unknown error finding finished event",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByNameType", mock.Anything, "test-task", app.EventTypeTaskFinished).
						Once().
						Return(app.Event{}, errors.New("something went wront"))
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test-task",
			},
			want:    app.CompletedTask{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewDurations(tt.fields.eventFinder)
			got, err := sut.FetchLastCompleted(tt.args.ctx, tt.args.taskName)
			tt.wantErr(t, err)
			assert.Equalf(t, tt.want, got, "FetchLastCompleted(%v, %v)", tt.args.ctx, tt.args.taskName)
			tt.fields.eventFinder.AssertExpectations(t)
		})
	}
}
