package aggregate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/screwyprof/cqrs"
)

// EventApplier applies events for the registered appliers.
type EventApplier struct {
	appliers map[string]cqrs.EventApplierFunc
}

// NewEventApplier creates a new instance of EventApplier.
func NewEventApplier() *EventApplier {
	return &EventApplier{
		appliers: make(map[string]cqrs.EventApplierFunc),
	}
}

// RegisterAppliers registers all the event appliers found in the aggregate.
func (a *EventApplier) RegisterAppliers(aggregate cqrs.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, "On") {
			continue
		}

		a.RegisterApplier(method.Name, func(e cqrs.DomainEvent) {
			method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(e)})
		})
	}
}

// RegisterApplier registers an event applier for the given method.
func (a *EventApplier) RegisterApplier(method string, applier cqrs.EventApplierFunc) {
	a.appliers[method] = applier
}

// Apply implements cqrs.EventApplier interface.
func (a *EventApplier) Apply(events ...cqrs.DomainEvent) error {
	for _, e := range events {
		if err := a.apply(e); err != nil {
			return err
		}
	}

	return nil
}

func (a *EventApplier) apply(event cqrs.DomainEvent) error {
	applierID := "On" + event.EventType()

	applier, ok := a.appliers[applierID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrEventApplierNotFound, applierID)
	}

	applier(event)

	return nil
}
