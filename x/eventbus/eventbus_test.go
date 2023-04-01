package eventbus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	event "github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/x"
	"github.com/screwyprof/cqrs/x/eventbus"
	"github.com/screwyprof/cqrs/x/eventhandler/evnhndtest"
)

// ensure that EventBus implements cqrs.EventPublisher interface.
var _ x.EventPublisher = (*eventbus.InMemoryEventBus)(nil)

func TestNewInMemoryEventBus(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventbus.NewInMemoryEventBus() != nil)
	})
}

func TestInMemoryEventBus_Publish(t *testing.T) {
	t.Run("ItFailsIfItCannotHandleAnEvent", func(t *testing.T) {
		// arrange
		eventHandler := &evnhndtest.EventHandlerMock{
			Err: evnhndtest.ErrCannotHandleEvent,
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(event.SomethingHappened{}, event.SomethingElseHappened{})

		// assert
		assert.Equal(t, evnhndtest.ErrCannotHandleEvent, err)
	})

	t.Run("ItPublishesEvents", func(t *testing.T) {
		// arrange
		want := []cqrs.DomainEvent{event.SomethingHappened{}, event.SomethingElseHappened{}}
		eventHandler := &evnhndtest.EventHandlerMock{}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(event.SomethingHappened{}, event.SomethingElseHappened{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, eventHandler.Happened)
	})

	t.Run("ItHandlesOnlyMatchedEvents", func(t *testing.T) {
		// arrange
		want := []cqrs.DomainEvent{event.SomethingHappened{}}
		eventHandler := &evnhndtest.EventHandlerMock{
			Matcher: cqrs.MatchEvent("SomethingHappened"),
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish([]cqrs.DomainEvent{
			event.SomethingHappened{},
			event.SomethingElseHappened{},
		}...)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, eventHandler.Happened)
	})
}
