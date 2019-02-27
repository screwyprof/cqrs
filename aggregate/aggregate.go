package aggregate

import (
	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

type Aggregate struct {
	id      uuid.UUID
	aggType string

	*eventApplier
	*commandHandler
	*changeRecorder
	*versionProvider
}

func NewAggregate(ID uuid.UUID, aggregateType string) *Aggregate {
	return &Aggregate{
		id:      ID,
		aggType: aggregateType,

		eventApplier:    newEventApplier(),
		commandHandler:  newCommandHandler(),
		changeRecorder:  newChangeRecorder(),
		versionProvider: newVersionProvider(),
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
	a.UpdateVersion(lastEvent.Version())

	return a.applyEvents(events...)
}

func (a *Aggregate) Apply(events ...cqrs.DomainEvent) error {
	if len(events) < 1 {
		return nil
	}

	a.recordChanges(events...)
	return a.applyEvents(events...)
}

func (a *Aggregate) recordChanges(events ...cqrs.DomainEvent) {
	for _, event := range events {
		event.SetAggregateID(a.id)
		event.SetVersion(a.nextEventVersion())

		a.recordChange(event)
	}
}

func (a *Aggregate) nextEventVersion() uint64 {
	return a.Version() + uint64(len(a.UncommittedChanges())) + 1
}
