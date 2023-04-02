package eventbus

import (
	"sync"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/x"
)

// InMemoryEventBus publishes events.
type InMemoryEventBus struct {
	eventHandlers   map[x.EventHandler]struct{}
	eventHandlersMu sync.RWMutex
}

// NewInMemoryEventBus creates a new instance of InMemoryEventBus.
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		eventHandlers: make(map[x.EventHandler]struct{}),
	}
}

// Register registers event handler.
func (b *InMemoryEventBus) Register(h x.EventHandler) {
	b.eventHandlersMu.Lock()
	defer b.eventHandlersMu.Unlock()

	b.eventHandlers[h] = struct{}{}
}

// Publish implements cqrs.EventPublisher interface.
func (b *InMemoryEventBus) Publish(events ...cqrs.DomainEvent) error {
	b.eventHandlersMu.RLock()
	defer b.eventHandlersMu.RUnlock()

	for h := range b.eventHandlers {
		if err := b.handleEvents(h, events...); err != nil {
			return err
		}
	}

	return nil
}

func (b *InMemoryEventBus) handleEvents(h x.EventHandler, events ...cqrs.DomainEvent) error {
	for _, e := range events {
		err := b.handleEventIfMatches(h.SubscribedTo(), h, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryEventBus) handleEventIfMatches(m cqrs.EventMatcher, h x.EventHandler, e cqrs.DomainEvent) error {
	if !m(e) {
		return nil
	}
	return h.Handle(e)
}
