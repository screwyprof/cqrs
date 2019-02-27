package aggregate

import (
	"fmt"
	"sync"

	"github.com/screwyprof/cqrs"
)

type commandHandler struct {
	handlers   map[string]cqrs.CommandHandlerFunc
	handlersMu sync.RWMutex
}

func newCommandHandler() *commandHandler {
	return &commandHandler{
		handlers: make(map[string]cqrs.CommandHandlerFunc),
	}
}

func (h *commandHandler) RegisterHandler(method string, handler cqrs.CommandHandlerFunc) {
	h.handlersMu.Lock()
	defer h.handlersMu.Unlock()
	h.handlers[method] = handler
}

func (h *commandHandler) Handle(c cqrs.Command) error {
	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	handler, ok := h.handlers[c.CommandType()]
	if !ok {
		return fmt.Errorf("handler for %s command is not found", c.CommandType())
	}

	return handler(c)
}
