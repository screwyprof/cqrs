package mock

import (
	"errors"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

var (
	// ErrAggregateStoreCannotLoadAggregate happens when aggregate store can't load aggregate.
	ErrAggregateStoreCannotLoadAggregate = errors.New("cannot load aggregate")
	// ErrAggregateStoreCannotStoreAggregate happens when aggregate store can't store aggregate.
	ErrAggregateStoreCannotStoreAggregate = errors.New("cannot store aggregate")
)

// AggregateStoreMock mocks event store.
type AggregateStoreMock struct {
	Loader func(aggregateID cqrs.Identifier, aggregateType string) (cqrs.AdvancedAggregate, error)
	Saver func(aggregate cqrs.AdvancedAggregate, events ...cqrs.DomainEvent) error
}

// Load implements cqrs.AggregateStore interface.
func (m *AggregateStoreMock) Load(
	aggregateID cqrs.Identifier, aggregateType string) (cqrs.AdvancedAggregate, error) {
	return m.Loader(aggregateID, aggregateType)
}

// StoreEventsFor implements cqrs.AggregateStore interface.
func (m *AggregateStoreMock) Store(aggregate cqrs.AdvancedAggregate, events ...cqrs.DomainEvent) error {
	return m.Saver(aggregate, events...)
}
