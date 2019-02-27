package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

// EventStore implements EventStore as an in memory structure.
type EventStore struct {
	committedEventStreams   map[uuid.UUID][]cqrs.DomainEvent
	committedEventStreamsMu sync.RWMutex

	uncommittedEventStreams   map[uuid.UUID][]cqrs.DomainEvent
	uncommittedEventStreamsMu sync.RWMutex
}

func NewEventStore() *EventStore {
	return &EventStore{
		committedEventStreams:   make(map[uuid.UUID][]cqrs.DomainEvent),
		uncommittedEventStreams: make(map[uuid.UUID][]cqrs.DomainEvent),
	}
}

func (s *EventStore) Store(eventProvider cqrs.EventProvider) error {
	fmt.Println("DB: Store")
	//if len(events) < 1 {
	//	return fmt.Errorf("no events given")
	//}

	version, err := s.loadEventProviderVersion(eventProvider)
	if err != nil {
		return err
	}

	if version != eventProvider.Version() {
		return fmt.Errorf("EventProvider has already been modified")
	}

	s.uncommittedEventStreamsMu.Lock()
	defer s.uncommittedEventStreamsMu.Unlock()

	s.uncommittedEventStreams[eventProvider.AggregateID()] = eventProvider.UncommittedChanges()

	eventProvider.UpdateVersion(eventProvider.Version() + uint64(len(eventProvider.UncommittedChanges())))
	// updateEventProviderVersion

	return nil
}

func (s *EventStore) loadEventProviderVersion(eventProvider cqrs.EventProvider) (uint64, error) {
	eventStream, err := s.loadEventStream(eventProvider.AggregateID())
	if err != nil {
		return 0, err
	}
	return uint64(len(eventStream)), nil
}

func (s *EventStore) LoadEventStream(aggregateID uuid.UUID) ([]cqrs.DomainEvent, error) {
	fmt.Println("DB: LoadEventStream")
	return s.loadEventStream(aggregateID)
}

func (s *EventStore) loadEventStream(aggregateID uuid.UUID) ([]cqrs.DomainEvent, error) {
	s.committedEventStreamsMu.RLock()
	defer s.committedEventStreamsMu.RUnlock()

	eventStream, ok := s.committedEventStreams[aggregateID]
	if !ok {
		return nil, nil
	}

	return eventStream, nil
}

func (s *EventStore) storeEvents(aggregateID uuid.UUID, events ...cqrs.DomainEvent) error {
	eventStream, err := s.loadEventStream(aggregateID)
	if err != nil {
		return err
	}

	eventStream = append(eventStream, events...)

	s.committedEventStreamsMu.Lock()
	s.committedEventStreams[aggregateID] = eventStream
	s.committedEventStreamsMu.Unlock()

	return nil
}

func (s *EventStore) printEvents(events ...cqrs.DomainEvent) {
	for _, event := range events {
		fmt.Printf("DB: Storing event: %s@%d of %s %+#v\n",
			event.EventID(), event.Version(), event.AggregateID().String(), event)
	}
}

func (s *EventStore) Commit() error {
	fmt.Println("DB: Commit")

	s.uncommittedEventStreamsMu.RLock()
	for aggregateID, uncommittedEventStream := range s.uncommittedEventStreams {
		if err := s.storeEvents(aggregateID, uncommittedEventStream...); err != nil {
			return err
		}
		s.printEvents(uncommittedEventStream...)
	}
	s.uncommittedEventStreamsMu.RUnlock()

	s.uncommittedEventStreamsMu.Lock()
	s.uncommittedEventStreams = make(map[uuid.UUID][]cqrs.DomainEvent)
	s.uncommittedEventStreamsMu.Unlock()

	return nil
}

func (s *EventStore) Rollback() error {
	fmt.Println("DB: Rollback")

	return nil
}

func (s *EventStore) BeginTransaction() {
	fmt.Println("DB: BeginTransaction")
}
