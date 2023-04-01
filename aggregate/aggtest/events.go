package aggtest

type SomethingHappened struct {
	Data string
}

func (c SomethingHappened) EventType() string {
	return "SomethingHappened"
}

type SomethingElseHappened struct{}

func (c SomethingElseHappened) EventType() string {
	return "SomethingElseHappened"
}
