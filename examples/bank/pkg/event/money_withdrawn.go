package event

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type MoneyWithdrawn struct {
	ID      domain.Identifier
	Amount  int64
	Balance int64
}

func (e MoneyWithdrawn) EventType() string {
	return "MoneyWithdrawn"
}
