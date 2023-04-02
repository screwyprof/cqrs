package x

import "github.com/screwyprof/cqrs"

// EventStore stores and loads events.
type EventStore interface {
	LoadEventsFor(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error)
	StoreEventsFor(aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent) error
}

// EventPublisher publishes events.
type EventPublisher interface {
	Publish(e ...cqrs.DomainEvent) error
}

// EventHandler handles events that were published though EventPublisher.
type EventHandler interface {
	SubscribedTo() cqrs.EventMatcher
	Handle(cqrs.DomainEvent) error
}

// EventHandlerFunc is a function that can be used as an event handler.
type EventHandlerFunc func(cqrs.DomainEvent) error

// AggregateStore loads and stores the aggregate.
type AggregateStore interface {
	Load(aggregateID cqrs.Identifier, aggregateType string) (cqrs.ESAggregate, error)
	Store(aggregate cqrs.ESAggregate, events ...cqrs.DomainEvent) error
}
