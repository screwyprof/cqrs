package aggregatestore

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

type Applier func(event cqrs.DomainEvent)
type Handler func(command cqrs.Command) error

type Aggregate struct {
	id uuid.UUID

	version      uint64
	eventVersion uint64

	appliedEvents []cqrs.DomainEvent

	appliers map[string]Applier
	handlers map[string]Handler
}

func NewAggregate(ID uuid.UUID) *Aggregate {
	return &Aggregate{
		id:       ID,
		appliers: make(map[string]Applier),
		handlers: make(map[string]Handler),
	}
}

func (a *Aggregate) AggregateID() uuid.UUID {
	return a.id
}

func (a *Aggregate) LoadFromHistory(events []cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	return a.applyChanges(events...)
}

func (a *Aggregate) UncommittedChanges() []cqrs.DomainEvent {
	return a.appliedEvents
}

func (a *Aggregate) MarkChangesAsCommitted() {
	a.appliedEvents = nil
}

func (a *Aggregate) Version() uint64 {
	return a.version
}

func (a *Aggregate) RegisterHandler(method string, handler Handler) {
	a.handlers[method] = handler
}

func (a *Aggregate) Handle(c cqrs.Command) error {
	commandType := reflect.TypeOf(c)

	handlerID := commandType.Name()
	handler, ok := a.handlers[handlerID]
	if !ok {
		return fmt.Errorf("handler for %s command is not found", handlerID)
	}

	return handler(c)
}

func (a *Aggregate) RegisterApplier(method string, applier Applier) {
	a.appliers[method] = applier
}

func (a *Aggregate) Apply(events ...cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	for _, e := range events {
		if err := a.apply(e); err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregate) apply(event cqrs.DomainEvent) error {
	event.SetAggregateID(a.id)
	event.SetVersion(a.nextEventVersion())

	if err := a.applyChange(event); err != nil {
		return err
	}

	a.appliedEvents = append(a.appliedEvents, event)
	a.version++

	return nil
}

func (a *Aggregate) applyChanges(events ...cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	var version uint64
	for _, e := range events {
		if err := a.applyChange(e); err != nil {
			return err
		}
		version++
	}

	lastEvent := events[len(events)-1]

	a.version = lastEvent.Version()
	a.eventVersion = a.version

	// maybe it's not necessary
	if a.version != version {
		return fmt.Errorf("last event version and calculated version aren't equal")
	}

	return nil
}

func (a *Aggregate) applyChange(event cqrs.DomainEvent) error {
	eventType := reflect.TypeOf(event)

	applierID := "On" + eventType.Name()
	applier, ok := a.appliers[applierID]
	if !ok {
		return fmt.Errorf("event handler for %s is not found", applierID)
	}
	applier(event)

	return nil
}

func (a *Aggregate) nextEventVersion() uint64 {
	a.eventVersion++
	return a.eventVersion
}
