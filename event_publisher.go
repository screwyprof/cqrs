package cqrs

// EventPublisher Publishes events.
type EventPublisher interface {
	Publish(...DomainEvent)
}
