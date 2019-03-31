package event

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type AccountOpened struct {
	ID     domain.Identifier
	Number string
}

func (e *AccountOpened) EventType() string {
	return "AccountOpened"
}
