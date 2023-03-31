package event

import "github.com/screwyprof/cqrs/examples/bank/domain"

// AccountOpened is an event which happens when an account is opened.
type AccountOpened struct {
	ID     domain.Identifier
	Number string
}

// EventType implements domain.Event interface.
func (e AccountOpened) EventType() string {
	return "AccountOpened"
}
