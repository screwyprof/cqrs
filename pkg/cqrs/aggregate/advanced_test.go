package aggregate_test

import (
	"testing"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	. "github.com/screwyprof/cqrs/pkg/cqrs/aggregate/testdata/fixture"
	. "github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

// ensure that Advanced implements cqrs.AdvancedAggregate interface.
var _ cqrs.AdvancedAggregate = (*aggregate.Advanced)(nil)

func TestNewBase(t *testing.T) {
	t.Run("ItPanicsIfThePureAggregateIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewAdvanced(nil, nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfCommandHandlerIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewAdvanced(NewTestAggregate(StringIdentifier("Test")), nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfEventApplierIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewAdvanced(
				NewTestAggregate(StringIdentifier("Test")),
				aggregate.NewCommandHandler(),
				nil,
			)
		}
		assert.Panic(t, factory)
	})
}

func TestBaseHandle(t *testing.T) {
	t.Run("ItUsesCustomCommandHandlerAndEventApplierWhenProvided", func(t *testing.T) {
		Test(t)(
			Given(createTestAggregateWithCustomCommandHandlerAndEventApplier()),
			When(MakeSomethingHappen{}),
			Then(SomethingHappened{}),
		)
	})

	t.Run("ItReturnsAnErrorIfTheHandlerIsNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyCommandHandler()),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrMakeSomethingHandlerNotFound),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier()),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrOnSomethingHappenedApplierNotFound),
		)
	})
}

func TestBaseVersion(t *testing.T) {
	t.Run("ItReturnsVersion", func(t *testing.T) {
		agg := createTestAggWithDefaultCommandHandlerAndEventApplier()

		assert.Equals(t, 0, agg.Version())
	})
}

func TestBaseApply(t *testing.T) {
	t.Run("ItAppliesEventsAndReturnsSomeBusinessError", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithDefaultCommandHandlerAndEventApplier(), SomethingHappened{}),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier(), SomethingHappened{}),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrOnSomethingHappenedApplierNotFound),
		)
	})

	t.Run("ItIncrementsVersion", func(t *testing.T) {
		agg := createTestAggWithEmptyCommandHandler()

		err := agg.Apply(SomethingHappened{})

		assert.Ok(t, err)
		assert.Equals(t, 1, agg.Version())
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.Advanced {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	handler := aggregate.NewCommandHandler()
	handler.RegisterHandlers(pureAgg)

	applier := aggregate.NewEventApplier()
	applier.RegisterAppliers(pureAgg)

	return aggregate.NewAdvanced(pureAgg, handler, applier)
}

func createTestAggregateWithCustomCommandHandlerAndEventApplier() *aggregate.Advanced {
	ID := StringIdentifier("TestAgg1")
	a := NewTestAggregate(ID)

	return aggregate.NewAdvanced(a, createCommandHandler(a), createEventApplier(a))
}

func createTestAggWithEmptyCommandHandler() *aggregate.Advanced {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	applier := aggregate.NewEventApplier()
	applier.RegisterAppliers(pureAgg)

	return aggregate.NewAdvanced(pureAgg, aggregate.NewCommandHandler(), applier)
}

func createTestAggWithEmptyEventApplier() *aggregate.Advanced {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	handler := aggregate.NewCommandHandler()
	handler.RegisterHandlers(pureAgg)

	return aggregate.NewAdvanced(pureAgg, handler, aggregate.NewEventApplier())
}

func createEventApplier(pureAgg *TestAggregate) *aggregate.EventApplier {
	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterApplier("OnSomethingHappened", func(e cqrs.DomainEvent) {
		pureAgg.OnSomethingHappened(e.(SomethingHappened))
	})
	return eventApplier
}

func createCommandHandler(pureAgg *TestAggregate) *aggregate.CommandHandler {
	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandler("MakeSomethingHappen", func(c cqrs.Command) ([]cqrs.DomainEvent, error) {
		return pureAgg.MakeSomethingHappen(c.(MakeSomethingHappen))
	})
	return commandHandler
}
