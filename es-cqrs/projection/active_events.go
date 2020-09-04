package projection

import (
	"es/event"
	"es/store"
	"time"
)

type ActiveEvent struct {
	ID        int
	Name      string
	StartTime time.Time
}

var ActiveEvents map[int]ActiveEvent = make(map[int]ActiveEvent)

func ProjectActiveEvents(evt store.SourceableEvent) error {
	switch e := evt.(type) {
	case event.EventStarted:
		ae := ActiveEvent{
			ID:        e.EventID,
			Name:      e.Name,
			StartTime: e.StartTime,
		}

		ActiveEvents[ae.ID] = ae
	case event.EventSuspended:
		delete(ActiveEvents, e.EventID)
	case event.EventClosed:
		delete(ActiveEvents, e.EventID)
	}

	return nil
}
