package store

import "github.com/screwyprof/cqrs/pkg/cqrs"

// AggregateStore loads and stores aggregates.
type AggregateStore struct {
	aggregateFactory cqrs.AggregateFactory
	eventStore       cqrs.EventStore
}

// NewStore creates a new instance of AggregateStore.
func NewStore(eventStore cqrs.EventStore, aggregateFactory cqrs.AggregateFactory) *AggregateStore {
	if eventStore == nil {
		panic("eventStore is required")
	}

	if aggregateFactory == nil {
		panic("aggregateFactory is required")
	}

	return &AggregateStore{
		eventStore:       eventStore,
		aggregateFactory: aggregateFactory,
	}
}

// Load implements cqrs.AggregateStore interface.
func (s *AggregateStore) Load(aggregateID cqrs.Identifier, aggregateType string) (cqrs.AdvancedAggregate, error) {
	loadedEvents, err := s.eventStore.LoadEventsFor(aggregateID)
	if err != nil {
		return nil, err
	}

	agg, err := s.aggregateFactory.CreateAggregate(aggregateType, aggregateID)
	if err != nil {
		return nil, err
	}

	err = agg.Apply(loadedEvents...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

// Store implements cqrs.AggregateStore interface.
func (s *AggregateStore) Store(agg cqrs.AdvancedAggregate, events ...cqrs.DomainEvent) error {
	return s.eventStore.StoreEventsFor(agg.AggregateID(), agg.Version(), events)
}
