package aggregate

import (
	"errors"
	"es/event"
	"es/eventstore"
)

type Event struct {
	ID   int
	Name string
	Type string
}

func LoadEvent(id int) (Event, error) {
	if events := eventstore.Get(id); events != nil {
		event := Event{}
		event.ApplyEvents(events)

		return event, nil
	}

	return Event{}, errors.New("Event not found")
}

func (e Event) ApplyEvents(events []eventstore.SourceableEvent) Event {
	for _, evt := range events {
		switch evt := evt.(type) {
		case event.CreateEventEvent:
			e.ID = evt.AggregateID
			e.Name = evt.Name
			e.Type = evt.Type
		}
	}
}
