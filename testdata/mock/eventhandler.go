package mock

import (
	"errors"

	cqrs2 "github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
)

var (
	ErrCannotHandleEvent    = errors.New("cannot handle event")
	ErrEventHandlerNotFound = errors.New("event handler for OnSomethingElseHappened event is not found")
)

type TestEventHandler struct {
	SomethingHappened string
}

func (h *TestEventHandler) OnSomethingHappened(e aggtest.SomethingHappened) error {
	h.SomethingHappened = e.Data
	return nil
}

func (h *TestEventHandler) OnSomethingElseHappened(e aggtest.SomethingElseHappened) error {
	return ErrCannotHandleEvent
}

func (h *TestEventHandler) SomeInvalidMethod() {
}

type EventHandlerMock struct {
	Err      error
	Matcher  cqrs2.EventMatcher
	Happened []cqrs2.DomainEvent
}

func (h *EventHandlerMock) SubscribedTo() cqrs2.EventMatcher {
	if h.Matcher != nil {
		return h.Matcher
	}
	return cqrs2.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")
}

func (h *EventHandlerMock) Handle(event cqrs2.DomainEvent) error {
	if h.Err != nil {
		return h.Err
	}
	switch e := event.(type) {
	case aggtest.SomethingHappened:
		h.OnSomethingHappened(e)
	case aggtest.SomethingElseHappened:
		h.OnSomethingElseHappened(e)
	}

	return nil
}

func (h *EventHandlerMock) OnSomethingHappened(e aggtest.SomethingHappened) {
	h.Happened = append(h.Happened, e)
}

func (h *EventHandlerMock) OnSomethingElseHappened(e aggtest.SomethingElseHappened) {
	h.Happened = append(h.Happened, e)
}
