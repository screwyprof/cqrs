package dispatcher_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	"github.com/screwyprof/cqrs/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/cqrs/pkg/cqrs/dispatcher/testdata/fixture"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
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
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreLoadErr(mock.ErrAggregateStoreCannotLoadAggregate),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrAggregateStoreCannotLoadAggregate),
		)
	})

	t.Run("ItFailsIfAggregateCannotHandleTheGivenCommand", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withLoadedEvents([]cqrs.DomainEvent{mock.SomethingHappened{}}),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItFailsIfItCannotStoreAggregate", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreSaveErr(mock.ErrAggregateStoreCannotStoreAggregate),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrAggregateStoreCannotStoreAggregate),
		)
	})

	t.Run("ItReturnsEvents", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		Test(t)(
			Given(createDispatcher(ID)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			Then(mock.SomethingHappened{}),
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

	pureAgg := mock.NewTestAggregate(id)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	agg := aggregate.NewAdvanced(pureAgg, commandHandler, eventApplier)
	if config.loadedEvents != nil {
		_ = agg.Apply(config.loadedEvents...)
	}

	return dispatcher.NewDispatcher(
		createAggregateStoreMock(agg, config.loadErr, config.storeErr),
	)
}

func createAggregateStoreMock(want cqrs.AdvancedAggregate, loadErr error, storeErr error) *mock.AggregateStoreMock {
	eventStore := &mock.AggregateStoreMock{
		Loader: func(aggregateID cqrs.Identifier, aggregateType string) (cqrs.AdvancedAggregate, error) {
			return want, loadErr
		},
		Saver: func(aggregate cqrs.AdvancedAggregate, events ...cqrs.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
