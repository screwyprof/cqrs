package account_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	. "github.com/screwyprof/cqrs/aggregate/aggtest/testdsl"
	"github.com/screwyprof/cqrs/examples/bank/domain"
	"github.com/screwyprof/cqrs/examples/bank/domain/account"
	"github.com/screwyprof/cqrs/examples/bank/domain/command"
	"github.com/screwyprof/cqrs/examples/bank/domain/event"
)

// ensure that the account aggregate implements cqrs.Aggregate interface.
var _ cqrs.Aggregate = (*account.Aggregate)(nil)

func TestAggregate(t *testing.T) {
	t.Run("panics if ID is not given", func(t *testing.T) {
		t.Parallel()

		factory := func() {
			account.NewAggregate(nil)
		}

		assert.Panics(t, factory)
	})

	t.Run("returns aggregate ID", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		agg := account.NewAggregate(ID)

		assert.Equal(t, ID, agg.AggregateID())
	})

	t.Run("returns aggregate type", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		agg := account.NewAggregate(ID)

		assert.Equal(t, "account.Aggregate", agg.AggregateType())
	})

	t.Run("opens an account", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()

		Test(t)(
			Given(createTestAggregate(ID)),
			When(command.OpenAccount{ID: ID, Number: number}),
			Then(event.AccountOpened{ID: ID, Number: number}),
		)
	})

	t.Run("deposits to an empty account", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()
		amount := faker.UnixTime()

		Test(t)(
			Given(createTestAggregate(ID), event.AccountOpened{ID: ID, Number: number}),
			When(command.DepositMoney{ID: ID, Amount: amount}),
			Then(event.MoneyDeposited{ID: ID, Amount: amount, Balance: amount}),
		)
	})

	t.Run("deposits to an account with initial funds", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
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

	t.Run("withdraws some funds", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()

		currentBalance := faker.UnixTime()
		amount := int64(100)
		newBalance := currentBalance - amount

		Test(t)(
			Given(createTestAggregate(ID),
				event.AccountOpened{ID: ID, Number: number},
				event.MoneyDeposited{ID: ID, Amount: currentBalance, Balance: currentBalance}),
			When(command.WithdrawMoney{ID: ID, Amount: amount}),
			Then(event.MoneyWithdrawn{ID: ID, Amount: amount, Balance: newBalance}),
		)
	})

	t.Run("cannot withdraw money if balance is not high enough", func(t *testing.T) {
		t.Parallel()

		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()
		amount := faker.UnixTime()

		Test(t)(
			Given(createTestAggregate(ID), event.AccountOpened{ID: ID, Number: number}),
			When(command.WithdrawMoney{ID: ID, Amount: amount}),
			ThenFailWith(account.ErrBalanceIsNotHighEnough),
		)
	})
}

func createTestAggregate(ID domain.Identifier) *aggregate.EventSourced {
	accAgg := account.NewAggregate(ID)

	return aggregate.FromAggregate(accAgg)
}
