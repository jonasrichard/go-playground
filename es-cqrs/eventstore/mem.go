package eventstore

type SourceableEvent interface {
	ID() int
}

var mem map[int][]SourceableEvent = make(map[int][]SourceableEvent)

func Store(event SourceableEvent) {
	events, ok := mem[event.ID()]
	if !ok {
		events = []SourceableEvent{event}
	} else {
		events = append(events, event)
	}

	mem[event.ID()] = events
}

func Get(id int) []SourceableEvent {
	events, ok := mem[id]
	if ok {
		return events
	}

	return nil
}
