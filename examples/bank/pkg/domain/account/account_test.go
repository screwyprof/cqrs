package account_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	. "github.com/screwyprof/cqrs/pkg/cqrs/aggregate/testdata/fixture"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/pkg/command"
	"github.com/screwyprof/cqrs/examples/bank/pkg/domain"
	"github.com/screwyprof/cqrs/examples/bank/pkg/domain/account"
	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
)

// ensure that game aggregate implements cqrs.Aggregate interface.
var _ cqrs.Aggregate = (*account.Aggregate)(nil)

func TestNewAggregate(t *testing.T) {
	t.Run("ItPanicsIfIDIsNotGiven", func(t *testing.T) {
		factory := func() {
			account.NewAggregate(nil)
		}
		assert.Panic(t, factory)
	})
}

func TestAggregateAggregateID(t *testing.T) {
	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		agg := account.NewAggregate(ID)

		assert.Equals(t, ID, agg.AggregateID())
	})
}

func TestAggregateAggregateType(t *testing.T) {
	t.Run("ItReturnsAggregateType", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		agg := account.NewAggregate(ID)

		assert.Equals(t, "account.Aggregate", agg.AggregateType())
	})
}

func TestAggregate(t *testing.T) {
	t.Run("ItOpensAnAccount", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()

		Test(t)(
			Given(createTestAggregate(ID)),
			When(command.OpenAccount{ID: ID, Number: number}),
			Then(event.AccountOpened{ID: ID, Number: number}),
		)
	})

	t.Run("ItDepositsAnEmptyAccount", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()
		amount := faker.UnixTime()

		Test(t)(
			Given(createTestAggregate(ID), event.AccountOpened{ID: ID, Number: number}),
			When(command.DepositMoney{ID: ID, Amount: amount}),
			Then(event.MoneyDeposited{ID: ID, Amount: amount, Balance: amount}),
		)
	})

	t.Run("ItDepositsAnAccountWithInitialFunds", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()

		currentBalance := faker.UnixTime()
		amount := faker.UnixTime()
		newBalance := amount + currentBalance

		Test(t)(
			Given(createTestAggregate(ID),
				event.AccountOpened{ID: ID, Number: number},
				event.MoneyDeposited{ID: ID, Amount: currentBalance, Balance: currentBalance},
			),
			When(command.DepositMoney{ID: ID, Amount: amount}),
			Then(event.MoneyDeposited{ID: ID, Amount: amount, Balance: newBalance}),
		)
	})
}

func createTestAggregate(ID domain.Identifier) *aggregate.Advanced {
	accAgg := account.NewAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(accAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(accAgg)

	return aggregate.NewAdvanced(accAgg, commandHandler, eventApplier)
}
