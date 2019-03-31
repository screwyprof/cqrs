package command

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type OpenAccount struct {
	ID     domain.Identifier
	Number string
}

func (c OpenAccount) AggregateID() domain.Identifier {
	return c.ID
}

func (c OpenAccount) AggregateType() string {
	return "account.Aggregate"
}

func (c OpenAccount) CommandType() string {
	return "OpenAccount"
}
