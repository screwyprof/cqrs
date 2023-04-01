package evnstoretest

import (
	"errors"

	"github.com/screwyprof/cqrs"
)

var (
	// ErrEventStoreCannotLoadEvents happens when event aggstore can't load events.
	ErrEventStoreCannotLoadEvents = errors.New("cannot load events")
	// ErrEventStoreCannotStoreEvents happens when event aggstore can't aggstore events.
	ErrEventStoreCannotStoreEvents = errors.New("cannot aggstore events")
)

// EventStoreMock mocks event aggstore.
type EventStoreMock struct {
	Loader func(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error)
	Saver  func(aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent) error
}

// LoadEventsFor implements cqrs.EventStore interface.
func (m *EventStoreMock) LoadEventsFor(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error) {
	return m.Loader(aggregateID)
}

// StoreEventsFor implements cqrs.EventStore interface.
func (m *EventStoreMock) StoreEventsFor(aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent) error {
	return m.Saver(aggregateID, version, events)
}
