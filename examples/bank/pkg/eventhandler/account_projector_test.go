package eventhandler_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	m "github.com/stretchr/testify/mock"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/internal/reporting"
	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
	eh "github.com/screwyprof/cqrs/examples/bank/pkg/eventhandler"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

func TestNewAccountDetailsProjector(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		projector := eh.NewAccountDetailsProjector(reporting.NewInMemoryAccountReporter())
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

		accountProjector := eh.NewAccountDetailsProjector(accountReporter)

		// act
		err := accountProjector.OnAccountOpened(event.AccountOpened{ID: ID, Number: number})

		// assert
		assert.Ok(t, err)
		accountReporter.AssertCalled(t, "Save", want)
	})
}

type accountReporterMock struct {
	m.Mock
}

func (r *accountReporterMock) AccountDetailsFor(ID report.Identifier) (*report.Account, error) {
	panic("implement me")
}

func (r *accountReporterMock) Save(account *report.Account) {
	r.Called(account)
}
