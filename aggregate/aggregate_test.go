package aggregate

import (
	"github.com/google/uuid"
	"github.com/screwyprof/cqrs"
	"reflect"
	"testing"
)

func TestNewAggregate(t *testing.T) {
	ID := uuid.New()
	agg := NewAggregate(ID, "TestAggregate")
	if agg == nil {
		t.Fatal("there should be an aggregate")
	}
	if agg.AggregateType() != "TestAggregate" {
		t.Error("the aggregate type should be correct: ", agg.AggregateType(), "TestAggregate")
	}
	if agg.AggregateID() != ID {
		t.Error("the aggregate ID should be correct: ", agg.AggregateID(), ID)
	}
	if agg.Version() != 0 {
		t.Error("the version should be 0:", agg.Version())
	}
}

//existingAggID := uuid.New()
//acc := account.Construct(existingAggID)
//
//accountOpened := event.NewAccountOpened(existingAggID, "ACC777")
//accountOpened.SetVersion(1)
//
//moneyDeposited := event.NewMoneyDeposited(existingAggID, 100, 100)
//moneyDeposited.SetVersion(2)
//
//moneyDeposited2 := event.NewMoneyDeposited(existingAggID, 50, 150)
//moneyDeposited2.SetVersion(3)
//
//err := acc.LoadFromHistory([]domain.DomainEvent{
//	accountOpened,
//	moneyDeposited,
//	moneyDeposited2,
//	//event.NewMoneyDeposited(existingAggID, 100, 100),
//	//event.NewMoneyDeposited(existingAggID, 50, 150),
//})
//failOnError(err)

//spew.Dump(acc.AggregateID())
//spew.Dump(acc.Version())
//spew.Dump(acc.UncommittedChanges())

//acc = account.Construct(uuid.New())
//err = acc.Handle(command.OpenAccount{Number:"ACC777"})
//failOnError(err)
//
//err = acc.Handle(command.DepositMoney{Amount:100})
//failOnError(err)

func TestAggregateEvents(t *testing.T) {
	ID := uuid.New()
	agg := NewTestAggregate(ID)

	event1 := NewTestEvent(ID, "event1")
	err := agg.Apply(event1)
	if err != nil {
		t.Errorf("no error expected, but got: %v", err)
	}
	if event1.EventType() != TestEventType {
		t.Error("the event type should be correct:", event1.EventType())
	}
	if !reflect.DeepEqual(event1.Content, "event1") {
		t.Error("the content should be correct:", event1.Content)
	}

	if event1.Version() != 1 {
		t.Error("the version should be 1:", event1.Version())
	}
	//if event1.AggregateType() != TestAggregateType {
	//	t.Error("the aggregate type should be correct:", event1.AggregateType())
	//}
	if event1.AggregateID() != ID {
		t.Error("the aggregate id should be correct:", event1.AggregateID())
	}
	//if event1.String() != "TestAggregateEvent@1" {
	//	t.Error("the string representation should be correct:", event1.String())
	//}
	events := agg.UncommittedChanges()
	if len(events) != 1 {
		t.Fatal("there should be one event stored:", len(events))
	}
	if events[0] != event1 {
		t.Error("the stored event should be correct:", events[0])
	}

	event2 := NewTestEvent(ID, "event1")
	err = agg.Apply(event2)
	if err != nil {
		t.Errorf("no error expected, but got: %v", err)
	}

	if event2.Version() != 2 {
		t.Error("the version should be 2:", event2.Version())
	}

	agg.MarkChangesAsCommitted()
	events = agg.UncommittedChanges()
	if len(events) != 0 {
		t.Error("there should be no events stored:", len(events))
	}

	event3 := NewTestEvent(ID, "event1")
	err = agg.Apply(event3)
	if err != nil {
		t.Errorf("no error expected, but got: %v", err)
	}

	if event3.Version() != 1 {
		t.Error("the version should be 1 after clearing uncommitted events (without applying any):", event3.Version())
	}

	agg = NewTestAggregate(uuid.New())
	event1 = NewTestEvent(ID, "event1")
	err = agg.Apply(event1)
	if err != nil {
		t.Errorf("no error expected, but got: %v", err)
	}

	event2 = NewTestEvent(ID, "event1")
	err = agg.Apply(event2)
	if err != nil {
		t.Errorf("no error expected, but got: %v", err)
	}

	events = agg.UncommittedChanges()
	if len(events) != 2 {
		t.Fatal("there should be 2 events stored:", len(events))
	}
	if events[0] != event1 {
		t.Error("the first stored event should be correct:", events[0])
	}
	if events[1] != event2 {
		t.Error("the second stored event should be correct:", events[0])
	}
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

type TestAggregateCommand struct {
	TestID  uuid.UUID
	Content string
}

func (t TestAggregateCommand) AggregateID() uuid.UUID { return t.TestID }
func (t TestAggregateCommand) AggregateType() string  { return TestAggregateType }
func (t TestAggregateCommand) CommandType() string    { return TestCommandType }

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
	*Aggregate
	event cqrs.DomainEvent
}

func NewTestAggregate(ID uuid.UUID) *TestAggregate {
	agg := &TestAggregate{
		Aggregate: NewAggregate(ID, TestAggregateType),
	}

	agg.RegisterApplier("OnTestEvent", func(e cqrs.DomainEvent) {
		agg.onTestEvent(e.(TestEvent))
	})

	return agg
}

func (a *TestAggregate) onTestEvent(e TestEvent) {
	a.event = e
}

//func (a *TestAggregate) Handle(cmd cqrs.Command) error {
//	return nil
//}

//func (a *TestAggregate) Apply(event cqrs.DomainEvent) error {
//	a.event = event
//	return nil
//}
