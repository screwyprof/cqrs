package aggregate

import (
	"sync"

	"github.com/screwyprof/cqrs"
)

type changeRecorder struct {
	uncommittedChanges   []cqrs.DomainEvent
	uncommittedChangesMu sync.RWMutex
}

func newChangeRecorder() *changeRecorder {
	return &changeRecorder{}
}

func (r *changeRecorder) UncommittedChanges() []cqrs.DomainEvent {
	r.uncommittedChangesMu.RLock()
	defer r.uncommittedChangesMu.RUnlock()
	return r.uncommittedChanges
}

func (r *changeRecorder) MarkChangesAsCommitted() {
	r.uncommittedChangesMu.Lock()
	defer r.uncommittedChangesMu.Unlock()
	r.uncommittedChanges = nil
}

func (r *changeRecorder) recordChange(event cqrs.DomainEvent) {
	r.uncommittedChangesMu.Lock()
	defer r.uncommittedChangesMu.Unlock()
	r.uncommittedChanges = append(r.uncommittedChanges, event)
}
