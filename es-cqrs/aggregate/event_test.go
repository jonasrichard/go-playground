package aggregate_test

import (
	"es/aggregate"
	"es/event"
	"es/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	events := []store.SourceableEvent{
		event.CreateEvent{EventID: 3, Name: "Man-Ars", Type: "Soccer"},
	}

	var e aggregate.Event

	e.ApplyEvents(events)

	assert.Equal(t, e.ID, 3)
	assert.Equal(t, e.Name, "Man-Ars")
	assert.Equal(t, e.Type, "Soccer")
}

func TestStartEvent(t *testing.T) {
	now := time.Now()
	events := []store.SourceableEvent{
		event.CreateEvent{EventID: 4, Name: "Man-Yrk", Type: "Soccer"},
		event.StartEvent{EventID: 4, StartTime: now},
	}

	var e aggregate.Event

	e.ApplyEvents(events)

	assert.Equal(t, e.ID, 4)
	//assert.Equal(t, e.StartTime, now)
}

func TestCreateMarket(t *testing.T) {
	events := []store.SourceableEvent{
		event.CreateEvent{EventID: 6, Name: "Deb-Man", Type: "Soccer"},
		event.CreateMarket{EventID: 6, MarketID: 2, Name: "Win-Draw-Win", Outcomes: []event.CreateMarketOutcome{
			{ID: 1, Name: "Home team", StartingPrice: 1.4},
			{ID: 2, Name: "Away team", StartingPrice: 5.4},
			{ID: 3, Name: "Draw", StartingPrice: 2.7},
		}},
	}

	var e aggregate.Event

	e.ApplyEvents(events)

	assert.Equal(t, e.ID, 6)
	assert.Equal(t, len(e.Markets), 1)
	assert.Equal(t, e.Markets[0].ID, 2)
	assert.Equal(t, e.Markets[0].Name, "Win-Draw-Win")
	assert.Equal(t, len(e.Markets[0].Outcomes), 3)
	assert.Contains(t, e.Markets[0].Outcomes, aggregate.Outcome{ID: 1, Name: "Home team", Price: 1.4})
	assert.Contains(t, e.Markets[0].Outcomes, aggregate.Outcome{ID: 2, Name: "Away team", Price: 5.4})
	assert.Contains(t, e.Markets[0].Outcomes, aggregate.Outcome{ID: 3, Name: "Draw", Price: 2.7})
}
