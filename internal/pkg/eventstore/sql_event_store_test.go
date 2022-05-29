package eventstore_test

import (
	"context"
	"database/sql"
	"github.com/danmurf/time-tracker/internal/app"
	"github.com/danmurf/time-tracker/internal/pkg/eventstore"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSQLEventStore_StoreFetchAll(t *testing.T) {
	event1 := app.Event{
		ID:        uuid.New(),
		Type:      app.EventTypeTaskStarted,
		TaskName:  "my-task-1",
		CreatedAt: time.Now().Add(-10 * time.Minute).Truncate(time.Second).UTC(),
	}
	event2 := app.Event{
		ID:        uuid.New(),
		Type:      app.EventTypeTaskFinished,
		TaskName:  "my-task-1",
		CreatedAt: time.Now().Add(-5 * time.Minute).Truncate(time.Second).UTC(),
	}
	event3 := app.Event{
		ID:        uuid.New(),
		Type:      app.EventTypeTaskStarted,
		TaskName:  "my-task-2",
		CreatedAt: time.Now().Add(-2 * time.Minute).Truncate(time.Second).UTC(),
	}
	type args struct {
		store []app.Event
	}
	tests := []struct {
		name string
		args args
		want []app.Event
	}{
		{
			name: "3 events",
			args: args{
				store: []app.Event{event1, event2, event3},
			},
			want: []app.Event{event1, event2, event3},
		},
		{
			name: "1 event",
			args: args{
				store: []app.Event{event1},
			},
			want: []app.Event{event1},
		},
		{
			name: "0 events",
			args: args{
				store: []app.Event{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			db := newMemorySqliteDB(t)
			defer db.Close()
			sut, err := eventstore.NewSQLEventStore(context.Background(), db)
			assert.NoError(t, err)

			for _, event := range tt.args.store {
				assert.NoError(t, sut.Store(ctx, event))
			}

			got, err := sut.FetchAll(ctx)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func newMemorySqliteDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("opening in memory sqlite database: %s", err)
		t.FailNow()
	}
	return db
}
