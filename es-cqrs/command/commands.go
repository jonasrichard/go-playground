package command

type CreateEventCommand struct {
	EventID int
	Name    string
	Type    string
}

type StartEventCommand struct {
	EventID int
}

type SuspendEventCommand struct {
	EventID int
}

type CloseEventCommand struct {
	EventID int
}

type CreateMarketOutcome struct {
	ID            int
	Name          string
	StartingPrice float64
}

type CreateMarketCommand struct {
	EventID  int
	MarketID int
	Name     string
	Outcomes []CreateMarketOutcome
}

type SuspendMarketCommand struct {
	EventID  int
	MarketID int
}

type UpdatePriceCommand struct {
	EventID   int
	MarketID  int
	OutcomeID int
	Price     float64
}
