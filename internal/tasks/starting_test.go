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

func TestStarter_Start(t *testing.T) {
	now := time.Now()
	nowFunc := func() time.Time {
		return now
	}
	id := uuid.New()
	uuidFunc := func() uuid.UUID {
		return id
	}

	type fields struct {
		eventStore *app_mocks.EventStore
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
			name: "successfully starts task",
			fields: fields{
				eventStore: func() *app_mocks.EventStore {
					m := &app_mocks.EventStore{}
					m.
						On("Store", mock.Anything, app.Event{
							ID:        id,
							Type:      app.EventTypeTaskStarted,
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
			name: "error storing event",
			fields: fields{
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
			sut := Starter{
				eventStore: tt.fields.eventStore,
				now:        nowFunc,
				newUUID:    uuidFunc,
			}
			tt.wantErr(t, sut.Start(tt.args.ctx, tt.args.taskName))
			tt.fields.eventStore.AssertExpectations(t)
		})
	}
}
