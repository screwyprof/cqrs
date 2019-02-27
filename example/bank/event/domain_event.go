package event

import "github.com/google/uuid"

type DomainEvent struct {
	ID    uuid.UUID
	AggID uuid.UUID

	eventType    string
	eventVersion uint64
}

func NewDomainEvent(aggID uuid.UUID, eventType string) *DomainEvent {
	return &DomainEvent{
		ID:        uuid.New(),
		AggID:     aggID,
		eventType: eventType,
	}
}

func (e *DomainEvent) EventID() uuid.UUID {
	return e.ID
}

func (e *DomainEvent) AggregateID() uuid.UUID {
	return e.AggID
}

func (e *DomainEvent) SetAggregateID(ID uuid.UUID) {
	e.AggID = ID
}

func (e *DomainEvent) EventType() string {
	return e.eventType
}

func (e *DomainEvent) SetVersion(version uint64) {
	e.eventVersion = version
}

func (e *DomainEvent) Version() uint64 {
	return e.eventVersion
}
