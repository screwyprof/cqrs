package aggregate

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/screwyprof/cqrs"
)

// CommandHandler registers and handles commands.
type CommandHandler struct {
	handlers   map[string]cqrs.CommandHandlerFunc
	handlersMu sync.RWMutex
}

// NewCommandHandler creates a new instance of CommandHandler.
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		handlers: make(map[string]cqrs.CommandHandlerFunc),
	}
}

// Handle implements cqrs.CommandHandler interface.
func (h *CommandHandler) Handle(c cqrs.Command) ([]cqrs.DomainEvent, error) {
	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	handler, ok := h.handlers[c.CommandType()]
	if !ok {
		return nil, fmt.Errorf("handler for %s command is not found", c.CommandType())
	}

	return handler(c)
}

// RegisterHandler registers a command handler for the given method.
func (h *CommandHandler) RegisterHandler(method string, handler cqrs.CommandHandlerFunc) {
	h.handlersMu.Lock()
	defer h.handlersMu.Unlock()
	h.handlers[method] = handler
}

// RegisterHandlers registers all the command handlers found in the aggregate.
func (h *CommandHandler) RegisterHandlers(aggregate cqrs.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !h.methodHasValidSignature(method) {
			continue
		}

		h.RegisterHandler(method.Name, func(c cqrs.Command) ([]cqrs.DomainEvent, error) {
			return h.invokeCommandHandler(method, aggregate, c)
		})
	}
}

func (h *CommandHandler) methodHasValidSignature(method reflect.Method) bool {
	if method.Type.NumIn() != 2 {
		return false
	}

	// ensure that the method has a cqrs.Command as a parameter.
	cmdIntfType := reflect.TypeOf((*cqrs.Command)(nil)).Elem()

	cmdType := method.Type.In(1)

	return cmdType.Implements(cmdIntfType)
}

func (h *CommandHandler) invokeCommandHandler(
	method reflect.Method, aggregate cqrs.Aggregate, c cqrs.Command,
) ([]cqrs.DomainEvent, error) {
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(c)})

	resErr := result[1].Interface()
	if resErr != nil {
		return nil, resErr.(error)
	}
	eventsIntf := result[0].Interface()

	events := h.convertDomainEvents(eventsIntf)
	return events, nil
}

func (h *CommandHandler) convertDomainEvents(eventsIntf interface{}) []cqrs.DomainEvent {
	eventsIntfs := h.interfaceSlice(eventsIntf)

	events := make([]cqrs.DomainEvent, 0, len(eventsIntfs))
	for _, eventIntf := range eventsIntfs {
		events = append(events, eventIntf.(cqrs.DomainEvent))
	}

	return events
}

func (h *CommandHandler) interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
