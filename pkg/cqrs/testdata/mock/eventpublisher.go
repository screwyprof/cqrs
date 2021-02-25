package mock

import (
	"errors"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

// ErrCannotPublishEvents happens when event publisher cannot publish the given events.
var ErrCannotPublishEvents = errors.New("cannot load aggregate")

// EventPublisherMock mocks event store.
type EventPublisherMock struct {
	Publisher func(e ...cqrs.DomainEvent) error
}

// Publish implements cqrs.EventPublisher interface.
func (m *EventPublisherMock) Publish(e ...cqrs.DomainEvent) error {
	return m.Publisher(e...)
}
