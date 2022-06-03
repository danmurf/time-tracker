package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/danmurf/time-tracker/internal/app"
	app_mocks "github.com/danmurf/time-tracker/internal/app/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestFinisher_Finish(t *testing.T) {
	now := time.Now()
	nowFunc := func() time.Time {
		return now
	}
	id := uuid.New()
	uuidFunc := func() uuid.UUID {
		return id
	}

	type fields struct {
		eventFinder *app_mocks.EventFinder
		eventStore  *app_mocks.EventStore
	}
	type args struct {
		ctx      context.Context
		taskName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successfully finishes task in progress",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByName", mock.Anything, "test").
						Once().
						Return(app.Event{
							ID:        uuid.UUID{},
							Type:      app.EventTypeTaskStarted,
							TaskName:  "test",
							CreatedAt: now.Add(-1 * time.Minute),
						}, nil)
					return m
				}(),
				eventStore: func() *app_mocks.EventStore {
					m := &app_mocks.EventStore{}
					m.
						On("Store", mock.Anything, app.Event{
							ID:        id,
							Type:      app.EventTypeTaskFinished,
							TaskName:  "test",
							CreatedAt: now,
						}).
						Once().
						Return(nil)
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test",
			},
			wantErr: assert.NoError,
		},
		{
			name: "unable to finish task already finished",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByName", mock.Anything, "test").
						Once().
						Return(app.Event{
							ID:        uuid.UUID{},
							Type:      app.EventTypeTaskFinished,
							TaskName:  "test",
							CreatedAt: now.Add(-1 * time.Minute),
						}, nil)
					return m
				}(),
				eventStore: &app_mocks.EventStore{},
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.True(t,
					errors.Is(err, app.ErrTaskNotStarted),
					fmt.Sprintf("want err [%s]; got [%s]", app.ErrTaskNotStarted, err),
				)
			},
		},
		{
			name: "unable to finish task never started",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByName", mock.Anything, "test").
						Once().
						Return(app.Event{}, app.ErrEventNotFound)
					return m
				}(),
				eventStore: &app_mocks.EventStore{},
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.True(t,
					errors.Is(err, app.ErrTaskNotStarted),
					fmt.Sprintf("want err [%s]; got [%s]", app.ErrTaskNotStarted, err),
				)
			},
		},
		{
			name: "unknown error finding latest event",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByName", mock.Anything, "test").
						Once().
						Return(app.Event{}, errors.New("something went wrong"))
					return m
				}(),
				eventStore: &app_mocks.EventStore{},
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test",
			},
			wantErr: assert.Error,
		},
		{
			name: "error storing event",
			fields: fields{
				eventFinder: func() *app_mocks.EventFinder {
					m := &app_mocks.EventFinder{}
					m.
						On("LatestByName", mock.Anything, "test").
						Once().
						Return(app.Event{
							ID:        uuid.UUID{},
							Type:      app.EventTypeTaskStarted,
							TaskName:  "test",
							CreatedAt: now.Add(-1 * time.Minute),
						}, nil)
					return m
				}(),
				eventStore: func() *app_mocks.EventStore {
					m := &app_mocks.EventStore{}
					m.
						On("Store", mock.Anything, mock.Anything).
						Once().
						Return(errors.New("something went wrong"))
					return m
				}(),
			},
			args: args{
				ctx:      context.Background(),
				taskName: "test",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := Finisher{
				eventStore:  tt.fields.eventStore,
				eventFinder: tt.fields.eventFinder,
				now:         nowFunc,
				newUUID:     uuidFunc,
			}
			tt.wantErr(t, sut.Finish(tt.args.ctx, tt.args.taskName))
			tt.fields.eventFinder.AssertExpectations(t)
			tt.fields.eventStore.AssertExpectations(t)
		})
	}
}
