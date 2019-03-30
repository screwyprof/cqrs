package eventstore_test

import (
	"testing"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventstore"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

// ensure that event store implements cqrs.EventStore interface.
var _ cqrs.EventStore = (*eventstore.InMemoryEventStore)(nil)

func TestNewInInMemoryEventStore(t *testing.T) {
	t.Run("ItCreatesEventStore", func(t *testing.T) {
		es := eventstore.NewInInMemoryEventStore()
		assert.True(t, es != nil)
	})
}

func TestInMemoryEventStoreLoadEventsFor(t *testing.T) {
	t.Run("ItLoadsEventsForTheGivenAggregate", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier("TestAgg")
		es := eventstore.NewInInMemoryEventStore()

		want := []cqrs.DomainEvent{mock.SomethingHappened{}}

		err := es.StoreEventsFor(ID, 0, want)
		assert.Ok(t, err)

		// act
		got, err := es.LoadEventsFor(ID)

		// assert
		assert.Ok(t, err)
		assert.Equals(t, want, got)
	})
}

func TestInMemoryEventStoreStoreEventsFor(t *testing.T) {
	t.Run("ItReturnsConcurrencyErrorIfVersionsAreNotTheSame", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier("TestAgg")
		es := eventstore.NewInInMemoryEventStore()

		// act
		err := es.StoreEventsFor(ID, 1, []cqrs.DomainEvent{mock.SomethingHappened{}})

		// assert
		assert.Equals(t, eventstore.ErrConcurrencyViolation, err)
	})
}
