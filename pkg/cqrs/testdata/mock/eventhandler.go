package mock

import (
	"errors"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

var (
	ErrCannotHandleEvent = errors.New("cannot handle event")
	ErrEventHandlerNotFound = errors.New("event handler for OnSomethingElseHappened event is not found")
)

type TestEventHandler struct {
	SomethingHappened string
}

func (h *TestEventHandler) OnSomethingHappened(e SomethingHappened) error {
	h.SomethingHappened = e.Data
	return nil
}

func (h *TestEventHandler) OnSomethingElseHappened(e SomethingElseHappened) error {
	return ErrCannotHandleEvent
}

func (h *TestEventHandler) SomeInvalidMethod() {

}

type EventHandlerMock struct {
	Err error
	Matcher cqrs.EventMatcher
	Happened []cqrs.DomainEvent
}

func (h *EventHandlerMock) SubscribedTo() cqrs.EventMatcher {
	if h.Matcher != nil {
		return h.Matcher
	}
	return cqrs.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")
}

func (h *EventHandlerMock) Handle(event cqrs.DomainEvent) error {
	if h.Err != nil {
		return h.Err
	}
	switch e := event.(type) {
	case SomethingHappened:
		h.OnSomethingHappened(e)
	case SomethingElseHappened:
		h.OnSomethingElseHappened(e)
	}

	return nil
}

func (h *EventHandlerMock) OnSomethingHappened(e SomethingHappened) {
	h.Happened = append(h.Happened, e)
}

func (h *EventHandlerMock) OnSomethingElseHappened(e SomethingElseHappened) {
	h.Happened = append(h.Happened, e)
}