package projection

type MarketCountry struct {
	EventID    int
	MarketID   int
	MarketName string
	Countries  []string
}

var MarketCountryView map[int]MarketCountry = make(map[int]MarketCountry)
