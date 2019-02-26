package bus

import (
	"github.com/screwyprof/cqrs"
)

type CommandHandler struct {
	repo cqrs.DomainRepository
}

func NewCommandHandler(repo cqrs.DomainRepository) cqrs.CommandHandler {
	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) Handle(c cqrs.Command) error {
	agg, err := h.repo.ByID(c.AggregateID(), c.AggregateType())
	if err != nil {
		return err
	}

	if err := agg.Handle(c); err != nil {
		return err
	}

	return err
}
