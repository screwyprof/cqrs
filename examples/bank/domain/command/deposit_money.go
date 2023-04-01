package command

import "github.com/screwyprof/cqrs/examples/bank/domain"

// DepositMoney is a command to credit an account.
type DepositMoney struct {
	ID     domain.Identifier
	Amount int64
}

// AggregateID implements cqrs.Command interface.
func (c DepositMoney) AggregateID() domain.Identifier {
	return c.ID
}

// AggregateType implements cqrs.Command interface.
func (c DepositMoney) AggregateType() string {
	return "account.Aggregate"
}

// CommandType implements cqrs.Command interface.
func (c DepositMoney) CommandType() string {
	return "DepositMoney"
}
