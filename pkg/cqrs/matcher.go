package cqrs

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(DomainEvent) bool

// MatchAny matches any event.
//func MatchAny() EventMatcher {
//	return func(e DomainEvent) bool {
//		return true
//	}
//}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(t string) EventMatcher {
	return func(e DomainEvent) bool {
		return e != nil && e.EventType() == t
	}
}

// MatchAnyEventOf matches if any of several matchers matches.
func MatchAnyEventOf(types ...string) EventMatcher {
	return func(e DomainEvent) bool {
		return matchAnyEvent(e, types...)
	}
}

func matchAnyEvent(e DomainEvent, types ...string) bool {
	for _, t := range types {
		if MatchEvent(t)(e) {
			return true
		}
	}
	return false
}
