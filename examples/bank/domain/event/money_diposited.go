package event

import "github.com/screwyprof/cqrs/examples/bank/domain"

// MoneyDeposited is an event which happens when an account is credited.
type MoneyDeposited struct {
	ID      domain.Identifier
	Amount  int64
	Balance int64
}

// EventType implements domain.Event interface.
func (e MoneyDeposited) EventType() string {
	return "MoneyDeposited"
}
