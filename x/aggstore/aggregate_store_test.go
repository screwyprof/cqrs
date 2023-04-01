package aggstore_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/x"
	"github.com/screwyprof/cqrs/x/aggstore"
	"github.com/screwyprof/cqrs/x/eventstore/evnstoretest"
)

// ensure that AggregateStore implements cqrs.AggregateStore interface.
var _ x.AggregateStore = (*aggstore.AggregateStore)(nil)

func TestNewStore(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggstore.NewStore(nil, nil)
		}
		assert.Panics(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggstore.NewStore(
				createEventStoreMock(nil, nil, nil),
				nil,
			)
		}
		assert.Panics(t, factory)
	})
}

func TestAggregateStoreLoad(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEventStoreLoadErr(evnstoretest.ErrEventStoreCannotLoadEvents))

		// act
		_, err := s.Load(ID, aggtest.TestAggregateType)

		// assert
		assert.Equal(t, evnstoretest.ErrEventStoreCannotLoadEvents, err)
	})

	t.Run("ItCannotCreateAggregate", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEmptyFactory())

		// act
		_, err := s.Load(ID, aggtest.TestAggregateType)

		// assert
		assert.Equal(t, aggtest.ErrAggIsNotRegistered, err)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(
			ID,
			withLoadedEvents([]cqrs.DomainEvent{aggtest.SomethingHappened{}}),
			withStaticEventApplier(),
		)

		// act
		_, err := s.Load(ID, aggtest.TestAggregateType)

		// assert
		assert.Equal(t, aggtest.ErrOnSomethingHappenedApplierNotFound, err)
	})

	t.Run("ItReturnsAggregate", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID)

		// act
		got, err := s.Load(ID, aggtest.TestAggregateType)

		// assert
		assert.NoError(t, err)
		assert.True(t, nil != got)
	})
}

func TestAggregateStoreStore(t *testing.T) {
	t.Run("ItFailsIfItCannotSafeEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEventStoreSaveErr(evnstoretest.ErrEventStoreCannotStoreEvents))
		agg := createAgg(ID)

		// act
		err := s.Store(agg, nil)

		// assert
		assert.Equal(t, evnstoretest.ErrEventStoreCannotStoreEvents, err)
	})
}

func createAgg(id cqrs.Identifier) *aggregate.EventSourced {
	pureAgg := aggtest.NewTestAggregate(id)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	return aggregate.New(pureAgg, commandHandler, eventApplier)
}

type aggregateStoreOptions struct {
	emptyFactory       bool
	staticEventApplier bool
	loadedEvents       []cqrs.DomainEvent

	loadErr  error
	storeErr error
}

type option func(*aggregateStoreOptions)

func withStaticEventApplier() option {
	return func(o *aggregateStoreOptions) {
		o.staticEventApplier = true
	}
}

func withEmptyFactory() option {
	return func(o *aggregateStoreOptions) {
		o.emptyFactory = true
	}
}

func withLoadedEvents(loadedEvents []cqrs.DomainEvent) option {
	return func(o *aggregateStoreOptions) {
		o.loadedEvents = loadedEvents
	}
}

func withEventStoreLoadErr(err error) option {
	return func(o *aggregateStoreOptions) {
		o.loadErr = err
	}
}

func withEventStoreSaveErr(err error) option {
	return func(o *aggregateStoreOptions) {
		o.storeErr = err
	}
}

func createAggregateStore(id cqrs.Identifier, opts ...option) *aggstore.AggregateStore {
	config := &aggregateStoreOptions{}
	for _, opt := range opts {
		opt(config)
	}

	pureAgg := aggtest.NewTestAggregate(id)

	applier := aggregate.NewEventApplier()
	if !config.staticEventApplier {
		applier.RegisterAppliers(pureAgg)
	}

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	agg := aggregate.New(pureAgg, commandHandler, applier)
	if config.loadedEvents != nil {
		_ = agg.Apply(config.loadedEvents...)
	}
	aggFactory := createAggFactory(agg, config.emptyFactory)
	eventStore := createEventStoreMock(config.loadedEvents, config.loadErr, config.storeErr)

	return aggstore.NewStore(eventStore, aggFactory)
}

func createAggFactory(agg *aggregate.EventSourced, empty bool) *aggregate.Factory {
	f := aggregate.NewFactory()
	if empty {
		return f
	}

	f.RegisterAggregate(func(ID cqrs.Identifier) cqrs.ESAggregate {
		return agg
	})

	return f
}

func createEventStoreMock(want []cqrs.DomainEvent, loadErr error, storeErr error) *evnstoretest.EventStoreMock {
	eventStore := &evnstoretest.EventStoreMock{
		Loader: func(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error) {
			return want, loadErr
		},
		Saver: func(aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
