package eventhandler_test

import (
	"errors"
	"testing"

	"github.com/bxcodec/faker/v4"
	m "github.com/stretchr/testify/mock"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventhandler"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
	eh "github.com/screwyprof/cqrs/examples/bank/pkg/eventhandler"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

func TestNewAccountDetailsProjector(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		projector := eh.NewAccountDetailsProjector(&accountReporterMock{})
		assert.True(t, projector != nil)
	})

	t.Run("ItPanicsIfAccountReporterIsNotGiven", func(t *testing.T) {
		factory := func() {
			eh.NewAccountDetailsProjector(nil)
		}
		assert.Panic(t, factory)
	})
}

func TestAccountDetailsProjector(t *testing.T) {
	t.Run("ItProjectsAccountOpenedEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()

		want := &report.Account{
			ID:     ID,
			Number: number,
		}

		accountReporter := &accountReporterMock{}
		accountReporter.On("Save", want)

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(event.AccountOpened{ID: ID, Number: number})

		// assert
		assert.Ok(t, err)
		accountReporter.AssertCalled(t, "Save", want)
	})

	t.Run("ItProjectsMoneyDepositedEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		amount := faker.UnixTime()
		balance := faker.UnixTime()

		accountReporter := createAccountReporterMock(eh.Ledger{
			ID:      ID,
			Action:  "deposit",
			Amount:  amount,
			Balance: balance,
		})

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(event.MoneyDeposited{ID: ID, Amount: amount, Balance: balance})

		// assert
		assert.Ok(t, err)
		accountReporter.AssertExpectations(t)
	})

	t.Run("ItReturnsAnErrorWhenItCannotProjectMoneyDepositedEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		want := errors.New("an error occurred")

		accountReporter := createAccountReporterMockWithError(ID, want)

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(
			event.MoneyDeposited{ID: ID, Amount: faker.UnixTime(), Balance: faker.UnixTime()})

		// assert
		assert.Equals(t, want, err)
		accountReporter.AssertExpectations(t)
	})

	t.Run("ItProjectsMoneyWithdrawnEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		amount := faker.UnixTime()
		balance := faker.UnixTime()

		accountReporter := createAccountReporterMock(eh.Ledger{
			ID:      ID,
			Action:  "withdraw",
			Amount:  amount,
			Balance: balance,
		})

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(event.MoneyWithdrawn{ID: ID, Amount: amount, Balance: balance})

		// assert
		assert.Ok(t, err)
		accountReporter.AssertExpectations(t)
	})

	t.Run("ItReturnsAnErrorWhenItCannotProjectMoneyWithdrawnEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		want := errors.New("an error occurred")

		accountReporter := createAccountReporterMockWithError(ID, want)

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(
			event.MoneyWithdrawn{ID: ID, Amount: faker.UnixTime(), Balance: faker.UnixTime()})

		// assert
		assert.Equals(t, want, err)
		accountReporter.AssertExpectations(t)
	})
}

func createAccountReporterMockWithError(ID report.Identifier, want error) *accountReporterMock {
	var accountReport *report.Account
	accountReporter := &accountReporterMock{}
	accountReporter.On("AccountDetailsFor", ID).Return(accountReport, want)
	return accountReporter
}

func createAccountReporterMock(l eh.Ledger) *accountReporterMock {
	number := faker.Word()
	want := &report.Account{
		ID:     l.ID,
		Number: number,
	}
	accountReporter := &accountReporterMock{}
	accountReporter.On("AccountDetailsFor", l.ID).Return(want, nil)

	want.Balance = l.Balance
	want.Ledgers = append(want.Ledgers, report.Ledger{Action: l.Action, Amount: l.Amount})
	accountReporter.On("Save", want)

	return accountReporter
}

type accountReporterMock struct {
	m.Mock
}

func (r *accountReporterMock) AccountDetailsFor(ID report.Identifier) (*report.Account, error) {
	args := r.Called(ID)
	return args.Get(0).(*report.Account), args.Error(1)
}

func (r *accountReporterMock) Save(account *report.Account) {
	r.Called(account)
}
