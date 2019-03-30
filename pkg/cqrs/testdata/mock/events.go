package mock

type SomethingHappened struct{}
func (c SomethingHappened) EventType() string {
	return "SomethingHappened"
}

type SomethingElseHappened struct{}
func (c SomethingElseHappened) EventType() string {
	return "SomethingElseHappened"
}