package aggregate_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/assert"
)

func TestNewAggregate(t *testing.T) {
	// arrange
	ID := uuid.New()

	// act
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	// assert
	assert.Assert(t, agg != nil, "there should be an aggregate")
}

func TestAggregateAggregateID(t *testing.T) {
	// arrange
	ID := uuid.New()

	// act
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	// assert
	assert.Assert(t, agg.AggregateID() == ID,
		"Invalid aggregate ID returned: got %v, expected %v ", agg.AggregateID(), ID)
}

func TestAggregateAggregateType(t *testing.T) {
	// arrange
	ID := uuid.New()

	// act
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	// assert
	assert.Assert(t, agg.AggregateType() == TestAggregateType,
		"Invalid aggregate type returned: got %v, expected %v ", agg.AggregateType(), TestAggregateType)
}

func TestAggregateLoadFromHistory_ValidEventStreamGiven_AggregateCorrectlyLoadedFromHistory(t *testing.T) {
	// arrange
	ID := uuid.New()

	event := NewTestEvent(ID, "event")
	event.SetVersion(1)

	eventStream := []cqrs.DomainEvent{
		event,
	}

	agg := NewTestAggregate(ID)

	// act
	err := agg.LoadFromHistory(eventStream)

	// assert
	assert.Ok(t, err)
	assert.Equals(t, uint64(1), agg.Version())
	assert.Equals(t, agg.event, eventStream[0])
}

func TestAggregateLoadFromHistory_EmptyEventStream_NothingChanged(t *testing.T) {
	// arrange
	agg := NewTestAggregate(uuid.New())

	// act
	err := agg.LoadFromHistory(nil)

	// assert
	assert.Ok(t, err)
}

func TestAggregateHandle_ValidCommandGiven_EventsApplied(t *testing.T) {
	// arrange
	ID := uuid.New()
	agg := NewTestAggregate(ID)

	c := TestCommand{
		TestID:  ID,
		Content: "test_data",
	}

	// act
	err := agg.Handle(c)
	expectedEvent := agg.event.(TestEvent)

	// assert
	assert.Ok(t, err)
	assert.Equals(t, ID, expectedEvent.AggregateID())
	assert.Equals(t, TestEventType, expectedEvent.EventType())
	assert.Equals(t, "test_data", expectedEvent.Content)
	assert.Equals(t, uint64(1), expectedEvent.Version())
}

func TestAggregateHandle_CommandHandlerNotRegistered_ErrorReturned(t *testing.T) {
	// arrange
	ID := uuid.New()
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	c := TestCommand{
		TestID:  ID,
		Content: "test_data",
	}

	expectedErrorText := "handler for TestCommand command is not found"

	// act
	err := agg.Handle(c)

	// assert
	assert.Assert(t, err.Error() == expectedErrorText, "expected error (%s), but got %s", expectedErrorText, err)
}

func TestAggregateApply_NoEventsGiven_NothingChanged(t *testing.T) {
	// arrange
	ID := uuid.New()
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	// act
	err := agg.Apply(nil)

	// assert
	assert.Ok(t, err)
}

func TestAggregateApply_EventHandlerNotRegistered_ErrorReturned(t *testing.T) {
	// arrange
	ID := uuid.New()
	agg := aggregate.NewAggregate(ID, TestAggregateType)

	expectedErrorText := "event handler for OnTestEvent is not found"

	// act
	err := agg.Apply(NewTestEvent(ID, "event"))

	// assert
	assert.Assert(t, err.Error() == expectedErrorText, "expected error (%s), but got %s", expectedErrorText, err)
}

func TestAggregateMarkChangesAsCommitted(t *testing.T) {
	// arrange
	ID := uuid.New()
	agg := NewTestAggregate(ID)

	err := agg.Apply(NewTestEvent(ID, "event"))
	assert.Ok(t, err)

	// act
	agg.MarkChangesAsCommitted()

	// assert
	assert.Equals(t, 0, len(agg.UncommittedChanges()))
}

func init() {
	cqrs.RegisterAggregate(func(id uuid.UUID) cqrs.ComplexAggregate {
		return NewTestAggregate(id)
	})
}

const (
	TestAggregateType = "TestAggregate"
	TestEventType     = "TestEvent"
	TestCommandType   = "TestCommand"
)

type TestCommand struct {
	TestID  uuid.UUID
	Content string
}

func (t TestCommand) AggregateID() uuid.UUID { return t.TestID }
func (t TestCommand) AggregateType() string  { return TestAggregateType }
func (t TestCommand) CommandType() string    { return TestCommandType }

type DomainEvent struct {
	ID    uuid.UUID
	AggID uuid.UUID

	eventType    string
	eventVersion uint64
}

func NewDomainEvent(aggID uuid.UUID, eventType string) *DomainEvent {
	return &DomainEvent{
		ID:        uuid.New(),
		AggID:     aggID,
		eventType: eventType,
	}
}

func (e *DomainEvent) EventID() uuid.UUID {
	return e.ID
}

func (e *DomainEvent) AggregateID() uuid.UUID {
	return e.AggID
}

func (e *DomainEvent) SetAggregateID(ID uuid.UUID) {
	e.AggID = ID
}

func (e *DomainEvent) EventType() string {
	return e.eventType
}

func (e *DomainEvent) SetVersion(version uint64) {
	e.eventVersion = version
}

func (e *DomainEvent) Version() uint64 {
	return e.eventVersion
}

type TestEvent struct {
	cqrs.DomainEvent

	Content string
}

func NewTestEvent(aggID uuid.UUID, content string) TestEvent {
	return TestEvent{
		DomainEvent: NewDomainEvent(aggID, TestEventType),
		Content:     content,
	}
}

type TestAggregate struct {
	*aggregate.Aggregate
	event cqrs.DomainEvent
}

func NewTestAggregate(ID uuid.UUID) *TestAggregate {
	agg := &TestAggregate{
		Aggregate: aggregate.NewAggregate(ID, TestAggregateType),
	}

	agg.RegisterHandler("TestCommand", func(c cqrs.Command) error {
		return agg.TestCommand(c.(TestCommand))
	})

	agg.RegisterApplier("OnTestEvent", func(e cqrs.DomainEvent) {
		agg.onTestEvent(e.(TestEvent))
	})

	return agg
}

func (a *TestAggregate) TestCommand(c TestCommand) error {
	return a.Apply(NewTestEvent(c.TestID, c.Content))
}

func (a *TestAggregate) onTestEvent(e TestEvent) {
	a.event = e
}
