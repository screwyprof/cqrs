package cqrs

type EventProvider interface {
	Aggregate

	Version() uint64
	//UpdateVersion(version uint64)

	LoadFromHistory(events []DomainEvent) error
	UncommittedChanges() []DomainEvent
	MarkChangesAsCommitted()
}
