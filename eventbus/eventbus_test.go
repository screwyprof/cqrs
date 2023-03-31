package eventbus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cqrs2 "github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/eventbus"
	"github.com/screwyprof/cqrs/testdata/mock"
)

// ensure that EventBus implements cqrs.EventPublisher interface.
var _ cqrs2.EventPublisher = (*eventbus.InMemoryEventBus)(nil)

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
		assert.Equal(t, mock.ErrCannotHandleEvent, err)
	})

	t.Run("ItPublishesEvents", func(t *testing.T) {
		// arrange
		want := []cqrs2.DomainEvent{mock.SomethingHappened{}, mock.SomethingElseHappened{}}
		eventHandler := &mock.EventHandlerMock{}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(mock.SomethingHappened{}, mock.SomethingElseHappened{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, eventHandler.Happened)
	})

	t.Run("ItHandlesOnlyMatchedEvents", func(t *testing.T) {
		// arrange
		want := []cqrs2.DomainEvent{mock.SomethingHappened{}}
		eventHandler := &mock.EventHandlerMock{
			Matcher: cqrs2.MatchEvent("SomethingHappened"),
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish([]cqrs2.DomainEvent{
			mock.SomethingHappened{},
			mock.SomethingElseHappened{},
		}...)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, eventHandler.Happened)
	})
}
