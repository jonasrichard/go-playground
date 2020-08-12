package aggregate

import (
	"errors"
	"es/event"
	"es/store"
	"fmt"
	"time"
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
	Price float32
}

type Market struct {
	ID       int
	Name     string
	Outcomes []Outcome
}

type Event struct {
	ID          int
	Name        string
	Type        string
	State       EventState
	StartTime   time.Time
	SuspendTime time.Time
	CloseTime   time.Time
	Markets     []Market
}

var ErrNotFound = errors.New("Event record not found")

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
		fmt.Printf("[Aggregate] Applying %T %v\n", evt, evt)

		switch evt := evt.(type) {
		case event.CreateEvent:
			e.ID = evt.EventID
			e.Name = evt.Name
			e.Type = evt.Type
			e.State = PreGame

		case event.StartEvent:
			// TODO put the validation into the command handler
			e.StartTime = evt.StartTime
			e.State = Started

		case event.SuspendEvent:
			// TODO cannot suspend event if it is not yet started
			e.SuspendTime = evt.SuspendTime

		case event.CloseEvent:
			// TODO cannot close event if it is not yet started or suspended
			e.CloseTime = evt.CloseTime

		case event.CreateMarket:
			outcomes := []Outcome{}

			for _, o := range evt.Outcomes {
				outcomes = append(outcomes, Outcome{ID: o.ID, Name: o.Name, Price: o.StartingPrice})
			}

			market := Market{
				ID:       evt.MarketID,
				Name:     evt.Name,
				Outcomes: outcomes,
			}

			if e.Markets == nil {
				e.Markets = []Market{market}
			} else {
				e.Markets = append(e.Markets, market)
			}

		case event.UpdatePrice:
			// look for outcome id and set the price
		}
	}

	fmt.Printf("[Aggregate] After applies %v\n", e)
}
