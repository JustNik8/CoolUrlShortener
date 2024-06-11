// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	domain "CoolUrlShortener/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// EventsServiceProducer is an autogenerated mock type for the EventsServiceProducer type
type EventsServiceProducer struct {
	mock.Mock
}

// ProduceEvent provides a mock function with given fields: event
func (_m *EventsServiceProducer) ProduceEvent(event domain.URLEvent) {
	_m.Called(event)
}

// NewEventsServiceProducer creates a new instance of EventsServiceProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventsServiceProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventsServiceProducer {
	mock := &EventsServiceProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}