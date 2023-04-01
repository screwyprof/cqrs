package store_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/testdata/mock"
	"github.com/screwyprof/cqrs/x/store"
)

// ensure that AggregateStore implements cqrs.AggregateStore interface.
var _ cqrs.AggregateStore = (*store.AggregateStore)(nil)

func TestNewStore(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(nil, nil)
		}
		assert.Panics(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(
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
		s := createAggregateStore(ID, withEventStoreLoadErr(mock.ErrEventStoreCannotLoadEvents))

		// act
		_, err := s.Load(ID, aggtest.TestAggregateType)

		// assert
		assert.Equal(t, mock.ErrEventStoreCannotLoadEvents, err)
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
		s := createAggregateStore(ID, withEventStoreSaveErr(mock.ErrEventStoreCannotStoreEvents))
		agg := createAgg(ID)

		// act
		err := s.Store(agg, nil)

		// assert
		assert.Equal(t, mock.ErrEventStoreCannotStoreEvents, err)
	})
}

func createAgg(id cqrs.Identifier) *aggregate.Advanced {
	pureAgg := aggtest.NewTestAggregate(id)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	return aggregate.NewAdvanced(pureAgg, commandHandler, eventApplier)
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

func createAggregateStore(id cqrs.Identifier, opts ...option) *store.AggregateStore {
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

	agg := aggregate.NewAdvanced(pureAgg, commandHandler, applier)
	if config.loadedEvents != nil {
		_ = agg.Apply(config.loadedEvents...)
	}
	aggFactory := createAggFactory(agg, config.emptyFactory)
	eventStore := createEventStoreMock(config.loadedEvents, config.loadErr, config.storeErr)

	return store.NewStore(eventStore, aggFactory)
}

func createAggFactory(agg *aggregate.Advanced, empty bool) *aggregate.Factory {
	f := aggregate.NewFactory()
	if empty {
		return f
	}
	f.RegisterAggregate(func(ID cqrs.Identifier) cqrs.AdvancedAggregate {
		return agg
	})

	return f
}

func createEventStoreMock(want []cqrs.DomainEvent, loadErr error, storeErr error) *mock.EventStoreMock {
	eventStore := &mock.EventStoreMock{
		Loader: func(aggregateID cqrs.Identifier) ([]cqrs.DomainEvent, error) {
			return want, loadErr
		},
		Saver: func(aggregateID cqrs.Identifier, version int, events []cqrs.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
