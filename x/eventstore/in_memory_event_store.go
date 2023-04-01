package eventstore

import (
	"errors"
	"sync"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/x"
)

// ErrConcurrencyViolation happens if aggregate has been modified concurrently.
var ErrConcurrencyViolation = errors.New("concurrency error: aggregate versions differ")

// InMemoryEventStore stores and loads events from memory.
type InMemoryEventStore struct {
	eventStreams   map[cqrs.Identifier][]cqrs.DomainEvent
	eventStreamsMu sync.RWMutex

	eventPublisher x.EventPublisher
}

// NewInInMemoryEventStore creates a new instance of InMemoryEventStore.
func NewInInMemoryEventStore(eventPublisher x.EventPublisher) *InMemoryEventStore {
	if eventPublisher == nil {
		panic("eventPublisher is required")
	}

	return &InMemoryEventStore{
		eventStreams:   make(map[cqrs.Identifier][]cqrs.DomainEvent),
		eventPublisher: eventPublisher,
	}
}

// LoadEventsFor loads events for the given aggregate.
func (s *InMemoryEventStore) LoadEventsFor(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error) {
	s.eventStreamsMu.RLock()
	defer s.eventStreamsMu.RUnlock()

	return s.eventStreams[aggregateID], nil
}

// StoreEventsFor saves evens of the given aggregate.
func (s *InMemoryEventStore) StoreEventsFor(
	aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent,
) error {
	previousEvents, _ := s.LoadEventsFor(aggregateID)
	if len(previousEvents) != version {
		return ErrConcurrencyViolation
	}

	s.eventStreamsMu.Lock()
	defer s.eventStreamsMu.Unlock()
	s.eventStreams[aggregateID] = events

	return s.eventPublisher.Publish(events...)
}
