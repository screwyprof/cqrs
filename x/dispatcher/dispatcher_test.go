package dispatcher_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/x/aggstore/aggstoretest"
	"github.com/screwyprof/cqrs/x/dispatcher"
	. "github.com/screwyprof/cqrs/x/dispatcher/testdsl"
)

// ensure that Dispatcher  implements cqrs.CommandHandler interface.
var _ cqrs.CommandHandler = (*dispatcher.Dispatcher)(nil)

func TestNewDispatcher(t *testing.T) {
	t.Run("ItPanicsIfAggregateStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil)
		}
		assert.Panics(t, factory)
	})
}

func TestNewDispatcherHandle(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadAggregate", func(t *testing.T) {
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreLoadErr(aggstoretest.ErrAggregateStoreCannotLoadAggregate),
			)),
			When(aggtest.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(aggstoretest.ErrAggregateStoreCannotLoadAggregate),
		)
	})

	t.Run("ItFailsIfAggregateCannotHandleTheGivenCommand", func(t *testing.T) {
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withLoadedEvents([]cqrs.DomainEvent{aggtest.SomethingHappened{}}),
			)),
			When(aggtest.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(aggtest.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItFailsIfItCannotStoreAggregate", func(t *testing.T) {
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreSaveErr(aggstoretest.ErrAggregateStoreCannotStoreAggregate),
			)),
			When(aggtest.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(aggstoretest.ErrAggregateStoreCannotStoreAggregate),
		)
	})

	t.Run("ItReturnsEvents", func(t *testing.T) {
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(ID)),
			When(aggtest.MakeSomethingHappen{AggID: ID}),
			Then(aggtest.SomethingHappened{}),
		)
	})
}

type dispatcherOptions struct {
	loadedEvents []cqrs.DomainEvent

	loadErr  error
	storeErr error
}

type option func(*dispatcherOptions)

func withLoadedEvents(loadedEvents []cqrs.DomainEvent) option {
	return func(o *dispatcherOptions) {
		o.loadedEvents = loadedEvents
	}
}

func withAggregateStoreLoadErr(err error) option {
	return func(o *dispatcherOptions) {
		o.loadErr = err
	}
}

func withAggregateStoreSaveErr(err error) option {
	return func(o *dispatcherOptions) {
		o.storeErr = err
	}
}

func createDispatcher(id cqrs.Identifier, opts ...option) *dispatcher.Dispatcher {
	config := &dispatcherOptions{}
	for _, opt := range opts {
		opt(config)
	}

	agg := aggtest.NewTestAggregate(id)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(agg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(agg)

	esAgg := aggregate.New(agg, commandHandler, eventApplier)
	if config.loadedEvents != nil {
		_ = esAgg.Apply(config.loadedEvents...)
	}

	return dispatcher.NewDispatcher(
		createAggregateStoreMock(esAgg, config.loadErr, config.storeErr),
	)
}

func createAggregateStoreMock(
	want cqrs.ESAggregate, loadErr error, storeErr error,
) *aggstoretest.AggregateStoreMock {
	eventStore := &aggstoretest.AggregateStoreMock{
		Loader: func(aggregateID cqrs.Identifier, aggregateType string) (cqrs.ESAggregate, error) {
			return want, loadErr
		},
		Saver: func(aggregate cqrs.ESAggregate, events ...cqrs.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
