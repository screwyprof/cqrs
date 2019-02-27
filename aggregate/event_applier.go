package aggregate

import (
	"fmt"
	"sync"

	"github.com/screwyprof/cqrs"
)

type Applier func(event cqrs.DomainEvent)

type eventApplier struct {
	appliers   map[string]Applier
	appliersMu sync.RWMutex
}

func newEventApplier() *eventApplier {
	return &eventApplier{
		appliers: make(map[string]Applier),
	}
}

func (a *eventApplier) RegisterApplier(method string, applier Applier) {
	a.appliersMu.Lock()
	defer a.appliersMu.Unlock()
	a.appliers[method] = applier
}

func (a *eventApplier) applyEvents(events ...cqrs.DomainEvent) error {
	for _, e := range events {
		if err := a.applyEvent(e); err != nil {
			return err
		}
	}

	return nil
}

func (a *eventApplier) applyEvent(event cqrs.DomainEvent) error {
	a.appliersMu.RLock()
	defer a.appliersMu.RUnlock()

	applierID := "On" + event.EventType()
	applier, ok := a.appliers[applierID]
	if !ok {
		return fmt.Errorf("event handler for %s is not found", applierID)
	}
	applier(event)

	return nil
}
