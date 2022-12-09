package store_test

import (
	"testing"

	"github.com/bxcodec/faker/v4"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	"github.com/screwyprof/cqrs/pkg/cqrs/store"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

// ensure that AggregateStore implements cqrs.AggregateStore interface.
var _ cqrs.AggregateStore = (*store.AggregateStore)(nil)

func TestNewStore(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(
				createEventStoreMock(nil, nil, nil),
				nil,
			)
		}
		assert.Panic(t, factory)
	})
}

func TestAggregateStoreLoad(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEventStoreLoadErr(mock.ErrEventStoreCannotLoadEvents))

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrEventStoreCannotLoadEvents, err)
	})

	t.Run("ItCannotCreateAggregate", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEmptyFactory())

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrAggIsNotRegistered, err)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(
			ID,
			withLoadedEvents([]cqrs.DomainEvent{mock.SomethingHappened{}}),
			withStaticEventApplier(),
		)

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrOnSomethingHappenedApplierNotFound, err)
	})

	t.Run("ItReturnsAggregate", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID)

		// act
		got, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Ok(t, err)
		assert.True(t, nil != got)
	})
}

func TestAggregateStoreStore(t *testing.T) {
	t.Run("ItFailsIfItCannotSafeEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		s := createAggregateStore(ID, withEventStoreSaveErr(mock.ErrEventStoreCannotStoreEvents))
		agg := createAgg(ID)

		// act
		err := s.Store(agg, nil)

		// assert
		assert.Equals(t, mock.ErrEventStoreCannotStoreEvents, err)
	})
}

func createAgg(ID cqrs.Identifier) *aggregate.Advanced {
	pureAgg := mock.NewTestAggregate(ID)

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

	loadErr      error
	storeErr     error
	publisherErr error
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

func createAggregateStore(ID cqrs.Identifier, opts ...option) *store.AggregateStore {
	config := &aggregateStoreOptions{}
	for _, opt := range opts {
		opt(config)
	}

	pureAgg := mock.NewTestAggregate(ID)

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
