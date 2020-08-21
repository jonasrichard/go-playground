package projection

import (
	"es/event"
	"es/store"
)

type MarketCountry struct {
	EventID    int
	MarketID   int
	MarketName string
	Countries  []string
}

var MarketCountryView map[int]MarketCountry = make(map[int]MarketCountry)

func ProjectMarketCountryView(evt store.SourceableEvent) {
	switch e := evt.(type) {
	case event.CreateMarket:
		countries := allowedCountries(e.Name)

		MarketCountryView[e.MarketID] = MarketCountry{
			EventID:    e.EventID,
			MarketID:   e.MarketID,
			MarketName: e.Name,
			Countries:  countries,
		}
	}
}

func allowedCountries(marketName string) []string {
	return []string{"uk", "de", "es"}
}
