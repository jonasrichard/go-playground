package aggregate

import (
	"errors"
	"es/event"
	"es/store"
	"fmt"
	"time"

	"github.com/fatih/color"
)

type EventState int

const (
	PreGame   EventState = 1
	Started   EventState = 2
	Suspended EventState = 3
	Closed    EventState = 4
)

type Outcome struct {
	ID    int
	Name  string
	Price float64
}

type Market struct {
	ID       int
	Name     string
	Outcomes map[int]Outcome
}

type Event struct {
	ID          int
	Name        string
	Type        string
	State       EventState
	StartTime   time.Time
	SuspendTime time.Time
	CloseTime   time.Time
	Markets     map[int]Market
}

var ErrNotFound = errors.New("Event record not found")
var Yellow = color.New(color.FgYellow)

func LoadEvent(id int) (*Event, error) {
	events := store.Get(id)

	if events != nil {
		event := &Event{}
		event.ApplyEvents(events)

		return event, nil
	}

	return nil, ErrNotFound
}

func (e *Event) ApplyEvents(events []store.SourceableEvent) {
	for _, evt := range events {
		Yellow.Print("[Aggregate]")
		fmt.Printf(" Applying %T %v\n", evt, evt)

		switch evt := evt.(type) {
		case event.EventCreated:
			e.ID = evt.EventID
			e.Name = evt.Name
			e.Type = evt.Type
			e.State = PreGame

		case event.EventUpdated:
			e.Name = evt.Name
			e.Type = evt.Type

		case event.EventStarted:
			e.StartTime = evt.StartTime
			e.State = Started

		case event.EventSuspended:
			e.SuspendTime = evt.SuspendTime
			e.State = Suspended

		case event.EventClosed:
			e.CloseTime = evt.CloseTime
			e.State = Closed

		case event.MarketCreated:
			outcomes := make(map[int]Outcome)

			for _, o := range evt.Outcomes {
				outcomes[o.ID] = Outcome{ID: o.ID, Name: o.Name, Price: o.StartingPrice}
			}

			market := Market{
				ID:       evt.MarketID,
				Name:     evt.Name,
				Outcomes: outcomes,
			}

			if e.Markets == nil {
				e.Markets = make(map[int]Market)
			}

			e.Markets[market.ID] = market

		case event.PriceUpdated:
			// look for outcome id and set the price
		}
	}

	Yellow.Print("[Aggregate]")
	fmt.Printf(" After applies %v\n", e)
}
