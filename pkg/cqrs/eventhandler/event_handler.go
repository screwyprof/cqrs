package eventhandler

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

// EventHandler handles events.
type EventHandler struct {
	handlers   map[string]cqrs.EventHandlerFunc
	handlersMu sync.RWMutex
}

// New creates new instance of New.
func New() *EventHandler {
	return &EventHandler{
		handlers: make(map[string]cqrs.EventHandlerFunc),
	}
}

// RegisterHandler registers an event handler for the given method.
func (h *EventHandler) RegisterHandler(method string, handler cqrs.EventHandlerFunc) {
	h.handlersMu.Lock()
	defer h.handlersMu.Unlock()
	h.handlers[method] = handler
}

// SubscribedTo implements cqrs.EventHandler interface.
func (h *EventHandler) SubscribedTo() cqrs.EventMatcher {
	var subscribedTo []string
	for m := range h.handlers {
		subscribedTo = append(subscribedTo, strings.TrimPrefix(m, "On"))
	}
	return cqrs.MatchAnyEventOf(subscribedTo...)
}

// Handle implements cqrs.EventHandler interface.
func (h *EventHandler) Handle(e cqrs.DomainEvent) error {
	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	handlerID := "On" + e.EventType()
	handler, ok := h.handlers[handlerID]
	if !ok {
		return fmt.Errorf("event handler for %s event is not found", handlerID)
	}

	return handler(e)
}

// RegisterHandlers registers all the event handlers found in .
func (h *EventHandler) RegisterHandlers(entity interface{}) {
	entityType := reflect.TypeOf(entity)
	for i := 0; i < entityType.NumMethod(); i++ {
		method := entityType.Method(i)
		h.registerHandlerDynamically(method, entity)
	}
}

func (h *EventHandler) registerHandlerDynamically(method reflect.Method, entity interface{}) {
	if !strings.HasPrefix(method.Name, "On") {
		return
	}

	h.RegisterHandler(method.Name, func(e cqrs.DomainEvent) error {
		return h.invokeEventHandler(method, entity, e)
	})
}

func (h *EventHandler) invokeEventHandler(method reflect.Method, entity interface{}, e cqrs.DomainEvent) error {
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(e)})
	resErr := result[0].Interface()
	if resErr != nil {
		return resErr.(error)
	}
	return nil
}
