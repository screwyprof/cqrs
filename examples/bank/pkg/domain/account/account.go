package account

import (
	"github.com/screwyprof/cqrs/examples/bank/pkg/command"
	"github.com/screwyprof/cqrs/examples/bank/pkg/domain"
	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
)

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
func (a *Aggregate) OpenAccount(c command.OpenAccount) ([]domain.DomainEvent, error) {
	return []domain.DomainEvent{event.AccountOpened{ID: c.ID, Number: c.Number}}, nil
}

// OnAccountOpened handles AccountOpened event.
func (a *Aggregate) OnAccountOpened(e event.AccountOpened) {
	a.number = e.Number
}

// DepositMoney Credits the account.
func (a *Aggregate) DepositMoney(c command.DepositMoney) ([]domain.DomainEvent, error) {
	balance := a.balance + c.Amount
	return []domain.DomainEvent{event.MoneyDeposited{ID: c.ID, Amount: c.Amount, Balance: balance}}, nil
}

// OnMoneyDeposited handles MoneyDeposited event.
func (a *Aggregate) OnMoneyDeposited(e event.MoneyDeposited) {
	a.balance = e.Balance
	//a.Ledgers = append(acc.Ledgers, Ledger{Action: "debit", Amount: e.Amount})
}
