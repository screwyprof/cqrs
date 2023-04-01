package evnhndtest

import (
	"errors"

	"github.com/screwyprof/cqrs"
	event "github.com/screwyprof/cqrs/aggregate/aggtest"
)

var (
	ErrCannotHandleEvent    = errors.New("cannot handle event")
	ErrEventHandlerNotFound = errors.New("event handler for OnSomethingElseHappened event is not found")
)

type TestEventHandler struct {
	SomethingHappened string
}

func (h *TestEventHandler) OnSomethingHappened(e event.SomethingHappened) error {
	h.SomethingHappened = e.Data
	return nil
}

func (h *TestEventHandler) OnSomethingElseHappened(_ event.SomethingElseHappened) error {
	return ErrCannotHandleEvent
}

func (h *TestEventHandler) SomeInvalidMethod() {
}

type EventHandlerMock struct {
	Err      error
	Matcher  cqrs.EventMatcher
	Happened []cqrs.DomainEvent
}

func (h *EventHandlerMock) SubscribedTo() cqrs.EventMatcher {
	if h.Matcher != nil {
		return h.Matcher
	}
	return cqrs.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")
}

func (h *EventHandlerMock) Handle(e cqrs.DomainEvent) error {
	if h.Err != nil {
		return h.Err
	}

	switch e := e.(type) {
	case event.SomethingHappened:
		h.OnSomethingHappened(e)
	case event.SomethingElseHappened:
		h.OnSomethingElseHappened(e)
	}

	return nil
}

func (h *EventHandlerMock) OnSomethingHappened(e event.SomethingHappened) {
	h.Happened = append(h.Happened, e)
}

func (h *EventHandlerMock) OnSomethingElseHappened(e event.SomethingElseHappened) {
	h.Happened = append(h.Happened, e)
}
