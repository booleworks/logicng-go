package maxsat

type EventMaxSatNewLowerBound struct {
	Bound int
}

func (EventMaxSatNewLowerBound) EventType() string {
	return "Max-SAT New Lower Bound"
}

type EventMaxSatNewUpperBound struct {
	Bound int
}

func (EventMaxSatNewUpperBound) EventType() string {
	return "Max-SAT New Upper Bound"
}
