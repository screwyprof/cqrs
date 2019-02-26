package cqrs

// EventHandler defines a standard interface for instances that wish to list for
// the occurrence of a specific event.
type EventHandler interface {
	// Handle allows an event to be "published" to interface implementations.
	// In the "real world", error handling would likely be implemented.
	Handle(DomainEvent)
}
