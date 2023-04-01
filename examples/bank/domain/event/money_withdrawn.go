package event

import "github.com/screwyprof/cqrs/examples/bank/domain"

// MoneyWithdrawn is an event which happens when an account is debited.
type MoneyWithdrawn struct {
	ID      domain.Identifier
	Amount  int64
	Balance int64
}

// EventType implements domain.Event interface.
func (e MoneyWithdrawn) EventType() string {
	return "MoneyWithdrawn"
}
