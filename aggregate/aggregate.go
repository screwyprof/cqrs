package aggregate

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

type Handler func(command cqrs.Command) error

type Aggregate struct {
	id      uuid.UUID
	aggType string
	version uint64

	uncommittedChanges []cqrs.DomainEvent

	eventApplier *eventApplier
	handlers     map[string]Handler
}

func NewAggregate(ID uuid.UUID, aggregateType string) *Aggregate {
	return &Aggregate{
		id:      ID,
		aggType: aggregateType,

		eventApplier: newEventApplier(),
		handlers:     make(map[string]Handler),
	}
}

func (a *Aggregate) AggregateID() uuid.UUID {
	return a.id
}

func (a *Aggregate) AggregateType() string {
	return a.aggType
}

func (a *Aggregate) LoadFromHistory(events []cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	lastEvent := events[len(events)-1]
	a.version = lastEvent.Version()

	return a.eventApplier.ApplyEvents(events...)
}

func (a *Aggregate) UncommittedChanges() []cqrs.DomainEvent {
	return a.uncommittedChanges
}

func (a *Aggregate) MarkChangesAsCommitted() {
	a.uncommittedChanges = nil
}

func (a *Aggregate) Version() uint64 {
	return a.version
}

func (a *Aggregate) UpdateVersion(version uint64) {
	a.version = version
}

func (a *Aggregate) RegisterHandler(method string, handler Handler) {
	a.handlers[method] = handler
}

func (a *Aggregate) Handle(c cqrs.Command) error {
	handler, ok := a.handlers[c.CommandType()]
	if !ok {
		return fmt.Errorf("handler for %s command is not found", c.CommandType())
	}

	return handler(c)
}

func (a *Aggregate) RegisterApplier(method string, applier Applier) {
	a.eventApplier.RegisterApplier(method, applier)
}

func (a *Aggregate) Apply(events ...cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	a.recordChanges(events...)
	return a.eventApplier.ApplyEvents(events...)
}

func (a *Aggregate) recordChanges(events ...cqrs.DomainEvent) {
	for _, event := range events {
		event.SetAggregateID(a.id)
		event.SetVersion(a.nextEventVersion())

		a.uncommittedChanges = append(a.uncommittedChanges, event)
	}
}

func (a *Aggregate) nextEventVersion() uint64 {
	return a.version + uint64(len(a.uncommittedChanges)) + 1
}
