// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	app "github.com/danmurf/time-tracker/internal/app"

	mock "github.com/stretchr/testify/mock"
)

// EventStore is an autogenerated mock type for the EventStore type
type EventStore struct {
	mock.Mock
}

// Store provides a mock function with given fields: ctx, event
func (_m *EventStore) Store(ctx context.Context, event app.Event) error {
	ret := _m.Called(ctx, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, app.Event) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewEventStoreT interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventStore creates a new instance of EventStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventStore(t NewEventStoreT) *EventStore {
	mock := &EventStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
