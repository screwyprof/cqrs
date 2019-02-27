package memory

import (
	"fmt"
	"github.com/screwyprof/cqrs"
)

type EventBus struct {
	// Using a map with an empty struct allows us to keep the handlers
	// unique while still keeping memory usage relatively low.
	handlers map[cqrs.EventHandler]struct{}
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[cqrs.EventHandler]struct{}),
	}
}

func (b *EventBus) Register(l cqrs.EventHandler) {
	b.handlers[l] = struct{}{}
}

func (b *EventBus) Deregister(l cqrs.EventHandler) {
	delete(b.handlers, l)
}

func (b *EventBus) Publish(events ...cqrs.DomainEvent) {
	for _, event := range events {
		b.publish(event)
	}
}

func (b *EventBus) publish(event cqrs.DomainEvent) {
	fmt.Printf("EventBus: Publishing event: %s@%d of %s %+#v\n",
		event.EventID(), event.Version(), event.AggregateID().String(), event)
	for h := range b.handlers {
		h.Handle(event)
	}
}
