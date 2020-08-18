package store

import (
	"fmt"

	"github.com/fatih/color"
)

type SourceableEvent interface {
	GetID() int
}

var mem map[int][]SourceableEvent = make(map[int][]SourceableEvent)
var magenta = color.New(color.FgMagenta)

func Store(event SourceableEvent) {
	id := event.GetID()
	events, ok := mem[id]

	if !ok {
		events = []SourceableEvent{event}
	} else {
		events = append(events, event)
	}

	magenta.Print("[Store]")
	fmt.Printf(" Store event %d -> %v\n", id, event)

	mem[id] = events
}

func Get(id int) []SourceableEvent {
	events, ok := mem[id]
	if ok {
		return events
	}

	return nil
}
