package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

// EventStore implements EventStore as an in memory structure.
type EventStore struct {
	eventStreams   map[uuid.UUID][]cqrs.DomainEvent
	eventStreamsMu sync.RWMutex
}

func NewEventStore() *EventStore {
	return &EventStore{
		eventStreams: make(map[uuid.UUID][]cqrs.DomainEvent),
	}
}

func (s *EventStore) Store(aggregateID uuid.UUID, version uint64, events []cqrs.DomainEvent) error {
	fmt.Println("DB: Store")
	if len(events) < 1 {
		return fmt.Errorf("no events given")
	}

	eventStream, err := s.getEventStream(aggregateID)
	if err != nil {
		return err
	}

	dbVersion := uint64(len(eventStream))

	currentVersion := dbVersion + uint64(len(events))
	if currentVersion != version {
		return fmt.Errorf("EventStream has already been modified (concurrently)")
	}

	s.eventStreamsMu.Lock()
	defer s.eventStreamsMu.Unlock()

	eventStream = append(eventStream, events...)
	s.eventStreams[aggregateID] = eventStream

	return nil
}

func (s *EventStore) LoadEventStream(aggregateID uuid.UUID) ([]cqrs.DomainEvent, error) {
	fmt.Println("DB: LoadEventStream")
	return s.getEventStream(aggregateID)
}

func (s *EventStore) getEventStream(aggregateID uuid.UUID) ([]cqrs.DomainEvent, error) {
	s.eventStreamsMu.RLock()
	defer s.eventStreamsMu.RUnlock()

	eventStream, ok := s.eventStreams[aggregateID]
	if !ok {
		return nil, nil
	}

	return eventStream, nil
}

func (s *EventStore) Commit() error {
	fmt.Println("DB: Commit")
	for _, eventStream := range s.eventStreams {
		for _, event := range eventStream {
			fmt.Printf("DB: Storing event: %s\n", event.EventType())
		}
	}
	return nil
}

func (s *EventStore) Rollback() error {
	fmt.Println("DB: Rollback")
	return nil
}

func (s *EventStore) BeginTransaction() {
	fmt.Println("DB: BeginTransaction")
}
