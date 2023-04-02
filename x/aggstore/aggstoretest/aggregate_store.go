package aggstoretest

import (
	"errors"

	"github.com/screwyprof/cqrs"
)

var (
	// ErrAggregateStoreCannotLoadAggregate happens when aggregate store can't load aggregate.
	ErrAggregateStoreCannotLoadAggregate = errors.New("cannot load aggregate")
	// ErrAggregateStoreCannotStoreAggregate happens when aggregate store can't save aggregate.
	ErrAggregateStoreCannotStoreAggregate = errors.New("cannot store aggregate")
)

// AggregateStoreMock mocks event store.
type AggregateStoreMock struct {
	Loader func(aggregateID cqrs.Identifier, aggregateType string) (cqrs.ESAggregate, error)
	Saver  func(aggregate cqrs.ESAggregate, events ...cqrs.DomainEvent) error
}

// Load implements cqrs.AggregateStore interface.
func (m *AggregateStoreMock) Load(
	aggregateID cqrs.Identifier, aggregateType string,
) (cqrs.ESAggregate, error) {
	return m.Loader(aggregateID, aggregateType)
}

// Store implements cqrs.AggregateStore interface.
func (m *AggregateStoreMock) Store(aggregate cqrs.ESAggregate, events ...cqrs.DomainEvent) error {
	return m.Saver(aggregate, events...)
}
