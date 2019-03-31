package command

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type WithdrawMoney struct {
	ID     domain.Identifier
	Amount int64
}

func (c WithdrawMoney) AggregateID() domain.Identifier {
	return c.ID
}

func (c WithdrawMoney) AggregateType() string {
	return "account.Aggregate"
}

func (c WithdrawMoney) CommandType() string {
	return "WithdrawMoney"
}
