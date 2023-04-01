package eventstore_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/x"
	"github.com/screwyprof/cqrs/x/eventbus/evnbustest"
	"github.com/screwyprof/cqrs/x/eventstore"
)

// ensure that event aggstore implements cqrs.EventStore interface.
var _ x.EventStore = (*eventstore.InMemoryEventStore)(nil)

func TestNewInInMemoryEventStore(t *testing.T) {
	t.Run("ItCreatesEventStore", func(t *testing.T) {
		es := eventstore.NewInInMemoryEventStore(createEventPublisherMock(nil))
		assert.True(t, es != nil)
	})

	t.Run("ItPanicsIfEventPublisherIsNotGiven", func(t *testing.T) {
		factory := func() {
			eventstore.NewInInMemoryEventStore(nil)
		}
		assert.Panics(t, factory)
	})
}

func TestInMemoryEventStoreLoadEventsFor(t *testing.T) {
	t.Run("ItLoadsEventsForTheGivenAggregate", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		es := eventstore.NewInInMemoryEventStore(createEventPublisherMock(nil))

		want := []cqrs.DomainEvent{aggtest.SomethingHappened{Data: faker.Word()}}

		err := es.StoreEventsFor(ID, 0, want)
		assert.NoError(t, err)

		// act
		got, err := es.LoadEventsFor(ID)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestInMemoryEventStoreStoreEventsFor(t *testing.T) {
	t.Run("ItReturnsConcurrencyErrorIfVersionsAreNotTheSame", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		es := eventstore.NewInInMemoryEventStore(createEventPublisherMock(nil))

		// act
		err := es.StoreEventsFor(ID, 1, []cqrs.DomainEvent{aggtest.SomethingHappened{}})

		// assert
		assert.Equal(t, eventstore.ErrConcurrencyViolation, err)
	})
}

func createEventPublisherMock(err error) *evnbustest.EventPublisherMock {
	eventPublisher := &evnbustest.EventPublisherMock{
		Publisher: func(e ...cqrs.DomainEvent) error {
			return err
		},
	}

	return eventPublisher
}
