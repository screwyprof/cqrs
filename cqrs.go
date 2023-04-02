package cqrs

import "fmt"

// Identifier represents an aggregate identifier.
type Identifier = fmt.Stringer

// Command is sent to the domain to change the state of an aggregate.
//
// Commands are named with a verb in the imperative mood, e.g., ConfirmOrder.
type Command interface {
	AggregateID() Identifier
	AggregateType() string
	CommandType() string
}

// CommandHandler is responsible for executing commands.
//
// It processes a command, produces relevant domain events.
//
// It returns a list of domain events on success
// It returns an error if the command cannot be executed.
type CommandHandler interface {
	Handle(c Command) ([]DomainEvent, error)
}

// CommandHandlerFunc is a function type that can be used as a command handler.
type CommandHandlerFunc func(Command) ([]DomainEvent, error)

// DomainEvent represents an event that has occurred in the domain.
//
// Events are named with a past-participle verb, e.g., OrderConfirmed.
type DomainEvent interface {
	EventType() string
}

// EventApplier is responsible for applying domain events to an aggregate.
type EventApplier interface {
	Apply(e ...DomainEvent) error
}

// EventApplierFunc is a function type that can be used as an event applier.
type EventApplierFunc func(DomainEvent)

// Aggregate represents a cluster of related objects that can be treated as a single unit.
//
// This basic interface is intended for simple aggregates that may not follow CQRS or event sourcing patterns.
type Aggregate interface {
	AggregateID() Identifier
	AggregateType() string
}

// Versionable indicates that an object can support different versions.
type Versionable interface {
	Version() int
}

// ESAggregate represents an aggregate that is designed with CQRS and event sourcing in mind.
//
// It extends the basic Aggregate interface and includes additional responsibilities such as
// command handling, event application, and versioning.
type ESAggregate interface {
	Aggregate
	Versionable
	CommandHandler
	EventApplier
}

// FactoryFn is a function type for an aggregate factory function.
type FactoryFn func(Identifier) ESAggregate

// AggregateFactory is responsible for creating aggregates.
// It registers aggregate factory functions and creates aggregates based on a given aggregate type and identifier.
type AggregateFactory interface {
	RegisterAggregate(aggregateType string, factory FactoryFn)
	CreateAggregate(aggregateType string, ID Identifier) (ESAggregate, error)
}
