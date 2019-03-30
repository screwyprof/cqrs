package eventbus_test

import (
	"testing"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventbus"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

// ensure that EventBus implements cqrs.EventPublisher interface.
var _ cqrs.EventPublisher = (*eventbus.InMemoryEventBus)(nil)

func TestNewInMemoryEventBus(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventbus.NewInMemoryEventBus() != nil)
	})
}

func TestInMemoryEventBus_Publish(t *testing.T) {
	t.Run("ItFailsIfItCannotHandleAnEvent", func(t *testing.T) {
		// arrange
		eventHandler := &mock.EventHandlerMock{
			Err: mock.ErrCannotHandleEvent,
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(mock.SomethingHappened{}, mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrCannotHandleEvent, err)
	})

	t.Run("ItPublishesEvents", func(t *testing.T) {
		// arrange
		want := []cqrs.DomainEvent{mock.SomethingHappened{}, mock.SomethingElseHappened{}}
		eventHandler := &mock.EventHandlerMock{}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(mock.SomethingHappened{}, mock.SomethingElseHappened{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, want, eventHandler.Happened)
	})

	t.Run("ItHandlesOnlyMatchedEvents", func(t *testing.T) {
		// arrange
		want := []cqrs.DomainEvent{mock.SomethingHappened{}}
		eventHandler := &mock.EventHandlerMock{
			Matcher: cqrs.MatchEvent("SomethingHappened"),
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish([]cqrs.DomainEvent{
			mock.SomethingHappened{},
			mock.SomethingElseHappened{},
		}...)

		// assert
		assert.Ok(t, err)
		assert.Equals(t, want, eventHandler.Happened)
	})
}
