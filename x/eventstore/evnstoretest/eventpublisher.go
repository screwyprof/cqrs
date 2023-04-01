package evnstoretest

import (
	"github.com/screwyprof/cqrs"
)

// EventPublisherMock mocks event store.
type EventPublisherMock struct {
	Publisher func(e ...cqrs.DomainEvent) error
}

// Publish implements cqrs.EventPublisher interface.
func (m *EventPublisherMock) Publish(e ...cqrs.DomainEvent) error {
	return m.Publisher(e...)
}
