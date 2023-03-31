package eventstore_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/eventstore"
	"github.com/screwyprof/cqrs/testdata/mock"
)

// ensure that event store implements cqrs.EventStore interface.
var _ cqrs.EventStore = (*eventstore.InMemoryEventStore)(nil)

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
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		es := eventstore.NewInInMemoryEventStore(createEventPublisherMock(nil))

		want := []cqrs.DomainEvent{mock.SomethingHappened{Data: faker.Word()}}

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
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		es := eventstore.NewInInMemoryEventStore(createEventPublisherMock(nil))

		// act
		err := es.StoreEventsFor(ID, 1, []cqrs.DomainEvent{mock.SomethingHappened{}})

		// assert
		assert.Equal(t, eventstore.ErrConcurrencyViolation, err)
	})
}

func createEventPublisherMock(err error) *mock.EventPublisherMock {
	eventPublisher := &mock.EventPublisherMock{
		Publisher: func(e ...cqrs.DomainEvent) error {
			return err
		},
	}

	return eventPublisher
}
