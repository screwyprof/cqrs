package cqrs

import "fmt"

// Identifier an object identifier.
type Identifier = fmt.Stringer

// Command is an object that is sent to the cqrs to change state.
//
// People request changes to the cqrs by sending commands.
// Command are named with a verb in the imperative mood, for example ConfirmOrder.
type Command interface {
	AggregateID() Identifier
	AggregateType() string
	CommandType() string
}

// CommandHandler executes commands.
//
// A command handler receives a command and brokers a result from the appropriate aggregate.
// "A result" is either a successful application of the command, or an error.
//
// Could be implemented like a method on the aggregate.
type CommandHandler interface {
	Handle(c Command) ([]DomainEvent, error)
}

// CommandHandlerFunc is a function that can be used as a command handler.
type CommandHandlerFunc func(Command) ([]DomainEvent, error)

// DomainEvent represents something that took place in the cqrs.
//
// Events are always named with a past-participle verb, such as OrderConfirmed.
type DomainEvent interface {
	EventType() string
}

// EventApplier applies the given events to an aggregate.
type EventApplier interface {
	Apply(e ...DomainEvent) error
}

// EventApplierFunc is a function that can be used as an event applier.
type EventApplierFunc func(DomainEvent)

// Aggregate is a cluster of cqrs objects that can be treated as a single unit.
//
// Every transaction is scoped to a single aggregate.
// The lifetimes of the components of an aggregate are bounded by
// the lifetime of the entire aggregate.
//
// Concretely, an aggregate will handle commands, apply events,
// and have a state model encapsulated within it that allows it to implement the required command validation,
// thus upholding the invariants (business rules) of the aggregate.
type Aggregate interface {
	AggregateID() Identifier
	AggregateType() string
}

type Versionable interface {
	Version() int
}

// AdvancedAggregate is an aggregate which handles commands
// and applies events after it automatically
type AdvancedAggregate interface {
	Aggregate
	Versionable
	CommandHandler
	EventApplier
}

// EventStore stores and loads events.
type EventStore interface {
	LoadEventsFor(aggregateID Identifier) ([]DomainEvent, error)
	StoreEventsFor(aggregateID Identifier, version int, events []DomainEvent) error
}

// FactoryFn aggregate factory function.
type FactoryFn func(Identifier) AdvancedAggregate

// AggregateFactory creates aggregates.
type AggregateFactory interface {
	RegisterAggregate(factory FactoryFn)
	CreateAggregate(aggregateType string, ID Identifier) (AdvancedAggregate, error)
}

// EventPublisher publishes events.
type EventPublisher interface {
	Publish(e ...DomainEvent) error
}

// EventHandler handles events that were published though EventPublisher.
type EventHandler interface {
	SubscribedTo() EventMatcher
	Handle(DomainEvent) error
}

// EventHandlerFunc is a function that can be used as an event handler.
type EventHandlerFunc func(DomainEvent) error

// AggregateStore loads and stores the aggregate.
type AggregateStore interface {
	Load(aggregateID Identifier, aggregateType string) (AdvancedAggregate, error)
	Store(aggregate AdvancedAggregate, events ...DomainEvent) error
}
