package projection_test

import (
	"es/event"
	"es/projection"
	"es/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivePrices(t *testing.T) {
	events := []store.SourceableEvent{
		event.EventCreated{EventID: 9, Name: "Man-Ars"},
		event.MarketCreated{EventID: 9, MarketID: 10, Name: "Win-Draw-Win", Outcomes: []event.CreateMarketOutcome{
			{ID: 11, Name: "Home team", StartingPrice: 6.5},
			{ID: 12, Name: "Draw", StartingPrice: 3.4},
			{ID: 13, Name: "Away team", StartingPrice: 2.1},
		}},
	}

	for _, evt := range events {
		err := projection.ProjectActiveEventPrice(evt)

		assert.Nil(t, err)
	}

	assert.NotNil(t, projection.ActiveEventPriceView[11])
	assert.Equal(t, projection.ActiveEventPriceView[11].EventID, 9)
	assert.Equal(t, projection.ActiveEventPriceView[11].MarketID, 10)
	assert.Equal(t, projection.ActiveEventPriceView[11].OutcomeID, 11)
	assert.Equal(t, projection.ActiveEventPriceView[11].Price, 6.5)
	assert.NotNil(t, projection.ActiveEventPriceView[11].ValidFrom)
}
