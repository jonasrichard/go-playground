package projection

import (
	"es/event"
	"es/store"
	"time"
)

type Func func(store.SourceableEvent) error

type ActivePrice struct {
	EventID   int
	MarketID  int
	OutcomeID int
	Price     float64
	ValidFrom time.Time
}

var ActiveEventPriceView map[int]ActivePrice = make(map[int]ActivePrice)

func ProjectActiveEventPrice(evt store.SourceableEvent) error {
	switch e := evt.(type) {
	case event.MarketCreated:
		for _, outcome := range e.Outcomes {
			price, ok := ActiveEventPriceView[outcome.ID]
			if !ok {
				price = ActivePrice{
					EventID:   e.EventID,
					MarketID:  e.MarketID,
					OutcomeID: outcome.ID,
					Price:     outcome.StartingPrice,
					ValidFrom: e.Timestamp,
				}
			} else {
				price.Price = outcome.StartingPrice
			}

			ActiveEventPriceView[outcome.ID] = price
		}

	case event.PriceUpdated:
		price, ok := ActiveEventPriceView[e.OutcomeID]
		if !ok {
			price = ActivePrice{
				EventID:   e.EventID,
				MarketID:  e.MarketID,
				OutcomeID: e.OutcomeID,
				Price:     e.Price,
				ValidFrom: e.Timestamp,
			}
		} else {
			price.Price = e.Price
		}

		ActiveEventPriceView[e.OutcomeID] = price
	}

	return nil
}
