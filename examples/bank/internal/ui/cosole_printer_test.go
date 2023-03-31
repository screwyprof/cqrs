package ui_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	m "github.com/stretchr/testify/mock"

	"github.com/screwyprof/cqrs/examples/bank/internal/ui"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
	"github.com/screwyprof/cqrs/testdata/mock"
)

func TestNewConsolePrinter(t *testing.T) {
	t.Run("ItPanicsIfWriterIsNotGiven", func(t *testing.T) {
		factory := func() {
			ui.NewConsolePrinter(nil, nil)
		}
		assert.Panics(t, factory)
	})

	t.Run("ItPanicsIfAccountReporterIsNotGiven", func(t *testing.T) {
		factory := func() {
			ui.NewConsolePrinter(&bytes.Buffer{}, nil)
		}
		assert.Panics(t, factory)
	})
}

func TestConsolePrinter_PrintAccountStatement(t *testing.T) {
	t.Run("ItPrintsDetailedAccountStatement", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())

		accReport := &report.Account{
			ID:      ID,
			Number:  faker.Word(),
			Balance: faker.RandomUnixTime(),
			Ledgers: []report.Ledger{
				{
					Action:  "deposit",
					Amount:  1000,
					Balance: 1000,
				},
				{
					Action:  "withdrawal",
					Amount:  100,
					Balance: 900,
				},
				{
					Action:  "deposit",
					Amount:  6000,
					Balance: 1500,
				},
			},
		}

		buf := &bytes.Buffer{}
		accountReporter := &accountReporterMock{}
		accountReporter.On("AccountDetailsFor", ID).Return(accReport, nil)

		printer := ui.NewConsolePrinter(buf, accountReporter)

		ledgers := bytes.Buffer{}
		for idx, ledger := range accReport.Ledgers {
			ledgers.WriteString(report.FormatLedger(idx+1, ledger))
		}

		want := fmt.Sprintf(
			"Account #%s:\n%s |%9s | %8s\n%s",
			accReport.Number,
			"#", "Amount", "Balance",
			ledgers.String(),
		)

		// act
		err := printer.PrintAccountStatement(ID)

		// assert
		assert.NoError(t, err)
		accountReporter.AssertExpectations(t)
		assert.Equal(t, want, buf.String())
	})

	t.Run("ItReturnsAnErrorIfItCannotPrintAccountStatement", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())

		want := fmt.Errorf("some error occured")

		var account *report.Account
		accountReporter := &accountReporterMock{}
		accountReporter.On("AccountDetailsFor", ID).Return(account, want)

		printer := ui.NewConsolePrinter(&bytes.Buffer{}, accountReporter)

		// act
		err := printer.PrintAccountStatement(ID)

		// assert
		assert.Equal(t, want, err)
	})
}

type accountReporterMock struct {
	m.Mock
}

func (r *accountReporterMock) AccountDetailsFor(ID report.Identifier) (*report.Account, error) {
	args := r.Called(ID)
	return args.Get(0).(*report.Account), args.Error(1)
}
