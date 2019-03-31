package command

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type DepositMoney struct {
	ID     domain.Identifier
	Amount int64
}

func (c DepositMoney) AggregateID() domain.Identifier {
	return c.ID
}

func (c DepositMoney) AggregateType() string {
	return "account.Aggregate"
}

func (c DepositMoney) CommandType() string {
	return "DepositMoney"
}
