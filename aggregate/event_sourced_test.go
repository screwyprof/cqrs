package aggregate_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	domain "github.com/screwyprof/cqrs/aggregate/aggtest"
	. "github.com/screwyprof/cqrs/aggregate/aggtest/testdsl"
)

// ensure that EventSourced implements cqrs.ESAggregate interface.
var _ cqrs.ESAggregate = (*aggregate.EventSourced)(nil)

func TestEventSourced(t *testing.T) {
	t.Parallel()

	t.Run("creating an event sourced aggregate", func(t *testing.T) {
		t.Parallel()

		t.Run("it panics if the aggregate is not provided", func(t *testing.T) {
			t.Parallel()

			factory := func() {
				aggregate.New(nil, nil, nil)
			}

			assert.Panics(t, factory)
		})

		t.Run("it panics if the command handler is not provided", func(t *testing.T) {
			t.Parallel()

			factory := func() {
				aggregate.New(domain.NewTestAggregate(domain.StringIdentifier(faker.UUIDHyphenated())), nil, nil)
			}

			assert.Panics(t, factory)
		})

		t.Run("it panics if the event applier is not provided", func(t *testing.T) {
			t.Parallel()

			factory := func() {
				aggregate.New(
					domain.NewTestAggregate(domain.StringIdentifier(faker.UUIDHyphenated())),
					aggregate.NewCommandHandler(),
					nil,
				)
			}

			assert.Panics(t, factory)
		})
	})

	t.Run("handling commands", func(t *testing.T) {
		t.Parallel()

		t.Run("it uses custom command handler and event applier when provided", func(t *testing.T) {
			t.Parallel()

			Test(t)(
				Given(createTestAggregateWithCustomCommandHandlerAndEventApplier()),
				When(domain.MakeSomethingHappen{}),
				Then(domain.SomethingHappened{}),
			)
		})

		t.Run("it returns an error if it cannot apply events", func(t *testing.T) {
			t.Parallel()

			Test(t)(
				Given(createTestAggWithEmptyEventApplier()),
				When(domain.MakeSomethingHappen{}),
				ThenFailWith(aggregate.ErrEventApplierNotFound),
			)
		})

		t.Run("it returns an error if the handler is not found", func(t *testing.T) {
			t.Parallel()

			Test(t)(
				Given(createTestAggWithEmptyCommandHandler()),
				When(domain.MakeSomethingHappen{}),
				ThenFailWith(aggregate.ErrCommandHandlerNotFound),
			)
		})

		t.Run("it returns an error if the command fails", func(t *testing.T) {
			t.Parallel()

			Test(t)(
				Given(createTestAggWithDefaultCommandHandlerAndEventApplier(), domain.SomethingHappened{}),
				When(domain.MakeSomethingHappen{}),
				ThenFailWith(domain.ErrItCanHappenOnceOnly),
			)
		})
	})

	t.Run("aggregate version", func(t *testing.T) {
		t.Parallel()

		t.Run("it returns the aggregate version", func(t *testing.T) {
			t.Parallel()

			agg := createTestAggWithDefaultCommandHandlerAndEventApplier()

			assert.Equal(t, 0, agg.Version())
		})
	})

	t.Run("applying events", func(t *testing.T) {
		t.Parallel()

		t.Run("it returns an error if the event appliers not found", func(t *testing.T) {
			agg := createTestAggWithEmptyEventApplier()

			err := agg.Apply(domain.SomethingHappened{})

			assert.ErrorIs(t, err, aggregate.ErrEventApplierNotFound)
		})

		t.Run("it increments the aggregate version", func(t *testing.T) {
			t.Parallel()

			agg := createTestAggWithEmptyCommandHandler()

			err := agg.Apply(domain.SomethingHappened{})

			assert.NoError(t, err)
			assert.Equal(t, 1, agg.Version())
		})
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.EventSourced {
	ID := domain.StringIdentifier(faker.UUIDHyphenated())
	agg := domain.NewTestAggregate(ID)

	handler := aggregate.NewCommandHandler()
	handler.RegisterHandlers(agg)

	applier := aggregate.NewEventApplier()
	applier.RegisterAppliers(agg)

	return aggregate.New(agg, handler, applier)
}

func createTestAggregateWithCustomCommandHandlerAndEventApplier() *aggregate.EventSourced {
	ID := domain.StringIdentifier(faker.UUIDHyphenated())
	agg := domain.NewTestAggregate(ID)

	return aggregate.New(agg, createCommandHandler(agg), createEventApplier(agg))
}

func createTestAggWithEmptyCommandHandler() *aggregate.EventSourced {
	ID := domain.StringIdentifier(faker.UUIDHyphenated())
	agg := domain.NewTestAggregate(ID)

	applier := aggregate.NewEventApplier()
	applier.RegisterAppliers(agg)

	return aggregate.New(agg, aggregate.NewCommandHandler(), applier)
}

func createTestAggWithEmptyEventApplier() *aggregate.EventSourced {
	ID := domain.StringIdentifier(faker.UUIDHyphenated())
	agg := domain.NewTestAggregate(ID)

	handler := aggregate.NewCommandHandler()
	handler.RegisterHandlers(agg)

	return aggregate.New(agg, handler, aggregate.NewEventApplier())
}

func createEventApplier(agg *domain.TestAggregate) *aggregate.EventApplier {
	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterApplier("OnSomethingHappened", func(e cqrs.DomainEvent) {
		agg.OnSomethingHappened(e.(domain.SomethingHappened)) //nolint:forcetypeassert
	})

	return eventApplier
}

func createCommandHandler(agg *domain.TestAggregate) *aggregate.CommandHandler {
	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(agg)

	return commandHandler
}
