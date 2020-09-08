package projection

import (
	"es/event"
	"es/store"
	"fmt"
	"time"

	"github.com/fatih/color"
)

type StartedEvent struct {
	ID        int
	Name      string
	StartTime time.Time
}

var Green = color.New(color.FgGreen).SprintFunc()

var AllEvents map[int]StartedEvent = make(map[int]StartedEvent)
var StartedEvents map[int]StartedEvent = make(map[int]StartedEvent)

func InitProjections() {
	ticker := time.Tick(10 * time.Second)
	go CheckStartTimePassed(ticker)
}

func CheckStartTimePassed(ticker <-chan time.Time) {
	for range ticker {
		now := time.Now()

		// Use mutex
		for id, evt := range AllEvents {
			if evt.StartTime.After(now) {
				StartedEvents[id] = evt

				fmt.Print(Green("[StartedEvent] "))
				fmt.Printf("Event %v is started\n", evt)
			}
		}
	}
}

func ProjectStartedEvents(evt store.SourceableEvent) error {
	switch e := evt.(type) {
	case event.EventCreated:
		se := StartedEvent{
			ID:        e.EventID,
			Name:      e.Name,
			StartTime: e.PlannedStartTime,
		}

		AllEvents[se.ID] = se
	}

	return nil
}
