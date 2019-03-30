package eventhandler_test

import (
	"testing"

	"github.com/screwyprof/cqrs/pkg/assert"

	"github.com/screwyprof/cqrs/examples/bank/internal/reporting"
	"github.com/screwyprof/cqrs/examples/bank/pkg/eventhandler"
)

func TestNewAccountDetailsProjector(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		projector := eventhandler.NewAccountDetailsProjector(reporting.NewInMemoryAccountReporter())
		assert.True(t, projector != nil)
	})
}
