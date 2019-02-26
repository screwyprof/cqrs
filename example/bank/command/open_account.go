package command

import (
	"github.com/google/uuid"
)

type OpenAccount struct {
	AggID   uuid.UUID
	AggType string

	Number string
}

func (c OpenAccount) AggregateID() uuid.UUID {
	return c.AggID
}

func (c OpenAccount) AggregateType() string {
	return "account.Account"
}
