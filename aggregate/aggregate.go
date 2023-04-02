// Package aggregate provides a base implementation for event sourced aggregates.
//
// To create an event sourced aggregate, a user must define their own domain
// aggregate that implements the cqrs.Aggregate interface. The user's domain
// aggregate should define its own identifier type and event type.
//
// Additionally, command handlers and event appliers should be defined within
// the user's domain aggregate. The command handlers handle commands and produce
// events, while the event appliers apply events to update the aggregate's state.
//
// Once the user's domain aggregate is defined, it can be transformed into
// a cqrs.ESAggregate using the aggregate.FromAggregate function. This function
// automatically registers the command handlers and event appliers defined in
// the user's domain aggregate.
//
// For a detailed example of how to use the aggregate package, please refer to
// the Example function in the example_test.go file.
//
// More examples can be found in the examples directory.
package aggregate

import (
	"errors"
)

var (
	// ErrCommandHandlerNotFound is returned when a command handler is not found.
	ErrCommandHandlerNotFound = errors.New("command handler not found")

	// ErrEventApplierNotFound is returned when an event applier is not found.
	ErrEventApplierNotFound = errors.New("event applier not found")

	// ErrAggregateNotRegistered is returned when an aggregate is not registered in the factory.
	ErrAggregateNotRegistered = errors.New("aggregate is not registered")
)
