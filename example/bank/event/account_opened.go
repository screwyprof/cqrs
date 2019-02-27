package event

import "github.com/google/uuid"

type AccountOpened struct {
	*DomainEvent

	Number string
}

func NewAccountOpened(aggID uuid.UUID, number string) AccountOpened {
	return AccountOpened{
		DomainEvent: NewDomainEvent(aggID, "AccountOpened"),
		Number:      number,
	}
}
