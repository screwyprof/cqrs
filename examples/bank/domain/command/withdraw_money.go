package command

import "github.com/screwyprof/cqrs/examples/bank/domain"

// WithdrawMoney is a command to debit an account.
type WithdrawMoney struct {
	ID     domain.Identifier
	Amount int64
}

// AggregateID implements cqrs.Command interface.
func (c WithdrawMoney) AggregateID() domain.Identifier {
	return c.ID
}

// AggregateType implements cqrs.Command interface.
func (c WithdrawMoney) AggregateType() string {
	return "account.Aggregate"
}

// CommandType implements cqrs.Command interface.
func (c WithdrawMoney) CommandType() string {
	return "WithdrawMoney"
}
