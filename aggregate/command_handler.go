package aggregate

import (
	"fmt"
	"reflect"

	"github.com/screwyprof/cqrs"
)

// CommandHandler registers and handles commands.
type CommandHandler struct {
	handlers map[string]cqrs.CommandHandlerFunc
}

// NewCommandHandler creates a new instance of CommandHandler.
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		handlers: make(map[string]cqrs.CommandHandlerFunc),
	}
}

// Handle implements cqrs.CommandHandler interface.
func (h *CommandHandler) Handle(c cqrs.Command) ([]cqrs.DomainEvent, error) {
	handler, ok := h.handlers[c.CommandType()]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrCommandHandlerNotFound, c.CommandType())
	}

	return handler(c)
}

// RegisterHandler registers a command handler for the given method.
func (h *CommandHandler) RegisterHandler(method string, handler cqrs.CommandHandlerFunc) {
	h.handlers[method] = handler
}

// RegisterHandlers registers all the command handlers found in the aggregate.
func (h *CommandHandler) RegisterHandlers(aggregate cqrs.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)

		if !h.isCommandHandler(method) {
			continue
		}

		h.RegisterHandler(method.Name, func(c cqrs.Command) ([]cqrs.DomainEvent, error) {
			return h.invokeCommandHandler(method, aggregate, c)
		})
	}
}

func (h *CommandHandler) isCommandHandler(method reflect.Method) bool {
	return h.commandHandlerHasExpectedInputs(method) && h.commandHandlerHasExpectedOutputs(method)
}

func (h *CommandHandler) commandHandlerHasExpectedInputs(method reflect.Method) bool {
	if method.Type.NumIn() != 2 {
		return false
	}

	cmdIntfType := reflect.TypeOf((*cqrs.Command)(nil)).Elem()
	cmdType := method.Type.In(1)

	return cmdType.Implements(cmdIntfType)
}

func (h *CommandHandler) commandHandlerHasExpectedOutputs(method reflect.Method) bool {
	return method.Type.NumOut() == 2 &&
		h.commandHandlerReturnsDomainEvents(method) &&
		h.commandHandlerReturnsAnError(method)
}

func (h *CommandHandler) commandHandlerReturnsDomainEvents(method reflect.Method) bool {
	eventSliceType := method.Type.Out(0)

	return eventSliceType.Kind() == reflect.Slice && h.isDomainEvent(eventSliceType.Elem())
}

func (h *CommandHandler) isDomainEvent(eventType reflect.Type) bool {
	method, ok := eventType.MethodByName("EventType")

	return ok && h.eventTypeMethodHasNoInputs(method) && h.eventTypeMethodReturnsString(method)
}

func (h *CommandHandler) commandHandlerReturnsAnError(method reflect.Method) bool {
	return method.Type.Out(1) == reflect.TypeOf((*error)(nil)).Elem()
}

func (h *CommandHandler) eventTypeMethodHasNoInputs(method reflect.Method) bool {
	return method.Type.NumIn() == 0
}

func (h *CommandHandler) eventTypeMethodReturnsString(method reflect.Method) bool {
	return method.Type.NumOut() == 1 && method.Type.Out(0) == reflect.TypeOf("")
}

func (h *CommandHandler) invokeCommandHandler(
	method reflect.Method, aggregate cqrs.Aggregate, c cqrs.Command,
) ([]cqrs.DomainEvent, error) {
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(c)})

	resErr := result[1].Interface()
	if resErr != nil {
		return nil, resErr.(error) //nolint:forcetypeassert
	}

	eventsIntf := result[0].Interface()
	events := h.convertDomainEvents(eventsIntf)

	return events, nil
}

func (h *CommandHandler) convertDomainEvents(eventsIntf interface{}) []cqrs.DomainEvent {
	eventsIntfs := h.interfaceSlice(eventsIntf)

	events := make([]cqrs.DomainEvent, 0, len(eventsIntfs))
	for _, eventIntf := range eventsIntfs {
		events = append(events, eventIntf.(cqrs.DomainEvent)) //nolint:forcetypeassert
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
