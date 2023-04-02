package aggregate

import (
	"github.com/screwyprof/cqrs"
)

// EventSourced is an aggregate that implements CQRS and Event Sourcing.
//
// It composes cqrs.Aggregate, cqrs.CommandHandler, and cqrs.EventApplier interfaces.
type EventSourced struct {
	cqrs.Aggregate
	version int

	commandHandler cqrs.CommandHandler
	eventApplier   cqrs.EventApplier
}

// New creates a new instance of EventSourced.
//
// It requires a basic cqrs.Aggregate, a cqrs.CommandHandler, and a cqrs.EventApplier as parameters.
// It returns a pointer to an EventSourced instance.
func New(aggregate cqrs.Aggregate, commandHandler cqrs.CommandHandler, eventApplier cqrs.EventApplier) *EventSourced {
	if aggregate == nil {
		panic("aggregate is required")
	}

	if commandHandler == nil {
		panic("commandHandler is required")
	}

	if eventApplier == nil {
		panic("eventApplier is required")
	}

	return &EventSourced{
		Aggregate:      aggregate,
		commandHandler: commandHandler,
		eventApplier:   eventApplier,
	}
}

// Version returns the current version of the aggregate.
//
// It implements the cqrs.Versionable interface.
func (b *EventSourced) Version() int {
	return b.version
}

// Handle processes the given command and produces relevant domain events.
//
// It handles this given command and applies the produced events.
//
// It implements the cqrs.CommandHandler interface.
func (b *EventSourced) Handle(c cqrs.Command) ([]cqrs.DomainEvent, error) {
	events, err := b.commandHandler.Handle(c)
	if err != nil {
		return nil, err
	}

	if applierErr := b.eventApplier.Apply(events...); applierErr != nil {
		return nil, applierErr
	}

	return events, nil
}

// Apply applies the given domain events to the aggregate.
//
// It applies the events and updates the aggregate version.
// It implements the cqrs.EventApplier interface.
func (b *EventSourced) Apply(e ...cqrs.DomainEvent) error {
	if err := b.eventApplier.Apply(e...); err != nil {
		return err
	}

	b.version += len(e)

	return nil
}
