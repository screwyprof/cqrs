package aggregate

import "github.com/screwyprof/cqrs/pkg/cqrs"

// Advanced implements an advanced aggregate root.
type Advanced struct {
	cqrs.Aggregate
	version int

	commandHandler cqrs.CommandHandler
	eventApplier   cqrs.EventApplier
}

// NewAdvanced creates a new instance of Advanced.
func NewAdvanced(pureAgg cqrs.Aggregate, commandHandler cqrs.CommandHandler, eventApplier cqrs.EventApplier) *Advanced {
	if pureAgg == nil {
		panic("pureAgg is required")
	}

	if commandHandler == nil {
		panic("commandHandler is required")
	}

	if eventApplier == nil {
		panic("eventApplier is required")
	}

	return &Advanced{
		Aggregate:      pureAgg,
		commandHandler: commandHandler,
		eventApplier:   eventApplier,
	}
}

// Version implements cqrs.Versionable interface.
func (b *Advanced) Version() int {
	return b.version
}

// Handle implements cqrs.CommandHandler.
func (b *Advanced) Handle(c cqrs.Command) ([]cqrs.DomainEvent, error) {
	events, err := b.commandHandler.Handle(c)
	if err != nil {
		return nil, err
	}

	if applierErr := b.eventApplier.Apply(events...); applierErr != nil {
		return nil, applierErr
	}

	return events, nil
}

// Apply implements cqrs.EventApplier interface.
func (b *Advanced) Apply(e ...cqrs.DomainEvent) error {
	if err := b.eventApplier.Apply(e...); err != nil {
		return err
	}
	b.version += len(e)
	return nil
}
