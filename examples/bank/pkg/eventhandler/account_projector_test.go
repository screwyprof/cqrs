package eventhandler_test

import (
	"errors"
	"testing"

	"github.com/bxcodec/faker/v3"
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
		number := faker.Word()
		amount := int(faker.UnixTime())
		balance := int(faker.UnixTime())

		want := &report.Account{
			ID:     ID,
			Number: number,
		}

		accountReporter := &accountReporterMock{}
		accountReporter.On("AccountDetailsFor", ID).Return(want, nil)

		want.Balance = balance
		want.Ledgers = append(want.Ledgers, report.Ledger{Action: "deposit", Amount: amount})

		accountReporter.On("Save", want)

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

		amount := int(faker.UnixTime())
		balance := int(faker.UnixTime())

		want := errors.New("an error occurred")

		var accountReport *report.Account

		accountReporter := &accountReporterMock{}
		accountReporter.On("AccountDetailsFor", ID).Return(accountReport, want)

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(event.MoneyDeposited{ID: ID, Amount: amount, Balance: balance})

		// assert
		assert.Equals(t, want, err)
		accountReporter.AssertExpectations(t)
	})

	t.Run("ItProjectsMoneyWithdrawnEvent", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		number := faker.Word()
		amount := int(faker.UnixTime())
		balance := int(faker.UnixTime())

		want := &report.Account{
			ID:     ID,
			Number: number,
		}

		accountReporter := &accountReporterMock{}
		accountReporter.On("AccountDetailsFor", ID).Return(want, nil)

		want.Balance = balance
		want.Ledgers = append(want.Ledgers, report.Ledger{Action: "withdraw", Amount: amount})

		accountReporter.On("Save", want)

		accountProjector := eventhandler.New()
		accountProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

		// act
		err := accountProjector.Handle(event.MoneyWithdrawn{ID: ID, Amount: amount, Balance: balance})

		// assert
		assert.Ok(t, err)
		accountReporter.AssertExpectations(t)
	})
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
