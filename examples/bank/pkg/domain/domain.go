package domain

import "fmt"

// Identifier an object identifier.
type Identifier = fmt.Stringer

//
// Event represents something that took place in the domain.
//
// Events are always named with a past-participle verb, such as OrderConfirmed.
type Event interface {
	EventType() string
}
