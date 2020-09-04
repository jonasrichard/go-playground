package event

import "time"

// TODO rename event to competition

type EventCreated struct {
	EventID   int
	Name      string
	Type      string
	Timestamp time.Time
}

type EventUpdated struct {
	EventID int
	Name    string
	Type    string
}

type EventStarted struct {
	EventID   int
	StartTime time.Time
	Name      string
}

type EventSuspended struct {
	EventID     int
	SuspendTime time.Time
}

type MarketResult struct {
	MarketID        int
	WinningOutcomes int
	LosingOutcomes  []int
}

type EventClosed struct {
	EventID   int
	CloseTime time.Time
	Result    []MarketResult
}

// helper struct for create market
type CreateMarketOutcome struct {
	ID            int
	Name          string
	StartingPrice float64
}

type MarketCreated struct {
	EventID   int
	MarketID  int
	Name      string
	Outcomes  []CreateMarketOutcome
	Timestamp time.Time
}

type PriceUpdated struct {
	EventID   int
	MarketID  int
	OutcomeID int
	Price     float64
	Timestamp time.Time
}

func (e EventCreated) GetID() int {
	return e.EventID
}

func (e EventUpdated) GetID() int {
	return e.EventID
}

func (e EventStarted) GetID() int {
	return e.EventID
}

func (e EventSuspended) GetID() int {
	return e.EventID
}

func (e EventClosed) GetID() int {
	return e.EventID
}

func (e MarketCreated) GetID() int {
	return e.EventID
}

func (e PriceUpdated) GetID() int {
	return e.EventID
}
