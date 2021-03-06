package reporting_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/internal/reporting"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

// ensure that account reporter implements report.AccountReporting interface.
var _ report.AccountReporting = (*reporting.InMemoryAccountReporter)(nil)

func TestNewInMemoryAccountReporter(t *testing.T) {
	t.Run("ItShouldCreateNewInstance", func(t *testing.T) {
		assert.True(t, reporting.NewInMemoryAccountReporter() != nil)
	})
}

func TestInMemoryAccountReporter(t *testing.T) {
	t.Run("ItShouldGetAccountDetailsForAGivenReportID", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		accountReporter := reporting.NewInMemoryAccountReporter()

		want := &report.Account{
			ID:      ID,
			Number:  faker.Word(),
			Balance: faker.UnixTime(),
			Ledgers: []report.Ledger{
				{
					Action: "debit",
					Amount: 100,
				},
				{
					Action: "withdraw",
					Amount: 50,
				},
			},
		}

		accountReporter.Save(want)

		got, err := accountReporter.AccountDetailsFor(ID)

		assert.Ok(t, err)
		assert.Equals(t, want, got)
	})

	t.Run("ItShouldReturnAccountNotFoundErrIfTheGivenAccountIsNotFound", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		accountReporter := reporting.NewInMemoryAccountReporter()

		_, err := accountReporter.AccountDetailsFor(ID)

		assert.Equals(t, reporting.ErrAccountNotFound, err)
	})
}
