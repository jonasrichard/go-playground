package event

type CreateEventEvent struct {
	AggregateID int
	EventID     int
	Name        string
	Type        string
}

func (e CreateEventEvent) ID() int {
	return e.AggregateID
}
