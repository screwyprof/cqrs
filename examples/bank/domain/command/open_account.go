package command

import "github.com/screwyprof/cqrs/examples/bank/domain"

// OpenAccount is a command to open an account.
type OpenAccount struct {
	ID     domain.Identifier
	Number string
}

// AggregateID implements cqrs.Command interface.
func (c OpenAccount) AggregateID() domain.Identifier {
	return c.ID
}

// AggregateType implements cqrs.Command interface.
func (c OpenAccount) AggregateType() string {
	return "account.Aggregate"
}

// CommandType implements cqrs.Command interface.
func (c OpenAccount) CommandType() string {
	return "OpenAccount"
}
