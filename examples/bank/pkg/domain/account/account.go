package account

import (
	"errors"

	"github.com/screwyprof/cqrs/examples/bank/pkg/command"
	"github.com/screwyprof/cqrs/examples/bank/pkg/domain"
	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
)

// ErrBalanceIsNotHighEnough happens when balance is not high enough.
var ErrBalanceIsNotHighEnough = errors.New("balance is not high enough")

// Aggregate handles operations with an account.
type Aggregate struct {
	id      domain.Identifier
	number  string
	balance int64
}

// NewAggregate creates a new instance of *Aggregate.
func NewAggregate(ID domain.Identifier) *Aggregate {
	if ID == nil {
		panic("ID required")
	}
	return &Aggregate{id: ID}
}

// AggregateID returns aggregate ID.
func (a *Aggregate) AggregateID() domain.Identifier {
	return a.id
}

// AggregateType return aggregate type.
func (a *Aggregate) AggregateType() string {
	return "account.Aggregate"
}

// OpenAccount opens a new account with a given number.
func (a *Aggregate) OpenAccount(c command.OpenAccount) ([]domain.Event, error) {
	return []domain.Event{event.AccountOpened{ID: c.ID, Number: c.Number}}, nil
}

// DepositMoney credits the account.
func (a *Aggregate) DepositMoney(c command.DepositMoney) ([]domain.Event, error) {
	balance := a.balance + c.Amount
	return []domain.Event{event.MoneyDeposited{ID: c.ID, Amount: c.Amount, Balance: balance}}, nil
}

// WithdrawMoney debits the account.
func (a *Aggregate) WithdrawMoney(c command.WithdrawMoney) ([]domain.Event, error) {
	balance := a.balance - c.Amount
	if balance <= 0 {
		return nil, ErrBalanceIsNotHighEnough
	}
	return []domain.Event{event.MoneyWithdrawn{ID: c.ID, Amount: c.Amount, Balance: balance}}, nil
}

// OnAccountOpened handles AccountOpened event.
func (a *Aggregate) OnAccountOpened(e event.AccountOpened) {
	a.number = e.Number
}

// OnMoneyDeposited handles MoneyDeposited event.
func (a *Aggregate) OnMoneyDeposited(e event.MoneyDeposited) {
	a.balance = e.Balance
	// a.Ledgers = append(acc.Ledgers, Ledger{Action: "debit", Amount: e.Amount})
}

// OnMoneyWithdrawn handles MoneyWithdrawn event.
func (a *Aggregate) OnMoneyWithdrawn(e event.MoneyWithdrawn) {
	a.balance = e.Balance
	// a.Ledgers = append(acc.Ledgers, Ledger{Action: "withdraw", Amount: e.Amount})
}
