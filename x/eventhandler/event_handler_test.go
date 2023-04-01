package eventhandler_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	event "github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/x/eventhandler"
	"github.com/screwyprof/cqrs/x/eventhandler/evnhndtest"
)

// ensure that event handler implements cqrs.EventHandler interface.
var _ cqrs.EventHandler = (*eventhandler.EventHandler)(nil)

func TestNew(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventhandler.New() != nil)
	})
}

func TestEventHandlerHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenEvent", func(t *testing.T) {
		// arrange
		eh := &evnhndtest.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandlers(eh)

		want := faker.Word()

		// act
		err := s.Handle(event.SomethingHappened{Data: want})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, eh.SomethingHappened)
	})

	t.Run("ItFailsIfEventHandlerIsNotRegistered", func(t *testing.T) {
		// arrange
		s := eventhandler.New()

		// act
		err := s.Handle(event.SomethingElseHappened{})

		// assert
		assert.Equal(t, evnhndtest.ErrEventHandlerNotFound, err)
	})

	t.Run("ItFailsIfEventHandlerReturnsAnError", func(t *testing.T) {
		// arrange
		eh := &evnhndtest.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandlers(eh)

		// act
		err := s.Handle(event.SomethingElseHappened{})

		// assert
		assert.Equal(t, evnhndtest.ErrCannotHandleEvent, err)
	})
}

func TestEventHandlerSubscribedTo(t *testing.T) {
	t.Run("ItReturnersTheEventsItSubscribedTo", func(t *testing.T) {
		// arrange
		eh := &evnhndtest.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandler("OnSomethingHappened", func(e cqrs.DomainEvent) error {
			return eh.OnSomethingHappened(e.(event.SomethingHappened)) //nolint:forcetypeassert
		})
		s.RegisterHandler("OnSomethingElseHappened", func(e cqrs.DomainEvent) error {
			return eh.OnSomethingElseHappened(e.(event.SomethingElseHappened)) //nolint:forcetypeassert
		})

		// act
		matcher := s.SubscribedTo()

		// assert
		assert.True(t, matcher(event.SomethingHappened{}))
		assert.True(t, matcher(event.SomethingElseHappened{}))
	})
}
