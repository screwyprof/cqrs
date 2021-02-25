package dispatcher

import "github.com/screwyprof/cqrs/pkg/cqrs"

// Dispatcher is a basic message dispatcher.
//
// It drives the overall command handling and event application/distribution process.
// It is suitable for a simple, single node application that can safely build its subscriber list
// at startup and keep it in memory.
// Depends on some kind of event storage mechanism.
type Dispatcher struct {
	store cqrs.AggregateStore
}

// NewDispatcher creates a new instance of Dispatcher.
func NewDispatcher(aggregateStore cqrs.AggregateStore) *Dispatcher {
	if aggregateStore == nil {
		panic("aggregateStore is required")
	}

	return &Dispatcher{
		store: aggregateStore,
	}
}

// Handle implements cqrs.CommandHandler interface.
func (d *Dispatcher) Handle(c cqrs.Command) ([]cqrs.DomainEvent, error) {
	agg, err := d.store.Load(c.AggregateID(), c.AggregateType())
	if err != nil {
		return nil, err
	}

	events, err := agg.Handle(c)
	if err != nil {
		return nil, err
	}

	if err = d.store.Store(agg, events...); err != nil {
		return nil, err
	}

	return events, nil
}
