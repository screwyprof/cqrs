package eventhandler_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventhandler"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
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
		eh := &mock.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandlers(eh)

		want := faker.Word()

		// act
		err := s.Handle(mock.SomethingHappened{Data: want})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, want, eh.SomethingHappened)
	})

	t.Run("ItFailsIfEventHandlerIsNotRegistered", func(t *testing.T) {
		// arrange
		s := eventhandler.New()

		// act
		err := s.Handle(mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrEventHandlerNotFound, err)
	})

	t.Run("ItFailsIfEventHandlerReturnsAnError", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandlers(eh)

		// act
		err := s.Handle(mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrCannotHandleEvent, err)
	})
}

func TestEventHandlerSubscribedTo(t *testing.T) {
	t.Run("ItReturnersTheEventsItSubscribedTo", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.New()
		s.RegisterHandler("OnSomethingHappened", func(e cqrs.DomainEvent) error {
			return eh.OnSomethingHappened(e.(mock.SomethingHappened))
		})
		s.RegisterHandler("OnSomethingElseHappened", func(e cqrs.DomainEvent) error {
			return eh.OnSomethingElseHappened(e.(mock.SomethingElseHappened))
		})

		// act
		matcher := s.SubscribedTo()

		// assert
		assert.True(t, matcher(mock.SomethingHappened{}))
		assert.True(t, matcher(mock.SomethingElseHappened{}))
	})
}
