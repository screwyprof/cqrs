package reporting_test

import (
	"testing"

	"github.com/screwyprof/cqrs/examples/bank/internal/reporting"
	"github.com/screwyprof/cqrs/pkg/assert"

	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

// ensure that game aggregate implements cqrs.Aggregate interface.
var _ report.AccountReporting = (*reporting.InMemoryAccountReporter)(nil)

func TestNewInMemoryAccountReporter(t *testing.T) {
	t.Run("ItShouldCreateNewInstance", func(t *testing.T) {
		assert.True(t, reporting.NewInMemoryAccountReporter() != nil)
	})
}
