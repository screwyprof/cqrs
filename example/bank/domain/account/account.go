package account

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"

	"github.com/screwyprof/cqrs/example/bank/command"
	"github.com/screwyprof/cqrs/example/bank/event"
)

type Account struct {
	*aggregate.Aggregate

	number  string
	balance int64
}

func Construct(ID uuid.UUID) *Account {
	acc := &Account{}
	acc.Aggregate = aggregate.NewAggregate(ID, "account.Account")

	acc.registerHandlers()
	acc.registerAppliers()

	return acc
}

func (a *Account) OpenAccount(c command.OpenAccount) error {
	return a.Apply(event.NewAccountOpened(c.AggregateID(), c.Number))
}

func (a *Account) DepositMoney(c command.DepositMoney) error {
	if err := a.guard(); err != nil {
		return err
	}

	balance := a.balance + c.Amount
	return a.Apply(event.NewMoneyDeposited(a.AggregateID(), c.Amount, balance))
}

func (a *Account) onAccountOpened(e event.AccountOpened) {
	a.number = e.Number
}

func (a *Account) onMoneyDeposited(e event.MoneyDeposited) {
	a.balance = e.Balance
}

// ToString Renders the Account as a string.
func (a *Account) ToString() string {
	return fmt.Sprintf("#%s: %d", a.number, a.balance)
}

func (a *Account) guard() error {
	if a.number == "" {
		return fmt.Errorf("account is not opened")
	}
	return nil
}

func (a *Account) registerHandlers() {
	a.RegisterHandler("OpenAccount", func(c cqrs.Command) error {
		return a.OpenAccount(c.(command.OpenAccount))
	})
	a.RegisterHandler("DepositMoney", func(c cqrs.Command) error {
		return a.DepositMoney(c.(command.DepositMoney))
	})
}

func (a *Account) registerAppliers() {
	a.RegisterApplier("OnAccountOpened", func(e cqrs.DomainEvent) {
		a.onAccountOpened(e.(event.AccountOpened))
	})
	a.RegisterApplier("OnMoneyDeposited", func(e cqrs.DomainEvent) {
		a.onMoneyDeposited(e.(event.MoneyDeposited))
	})
}
