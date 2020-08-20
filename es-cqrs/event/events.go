package event

import "time"

type CreateEvent struct {
	EventID   int
	Name      string
	Type      string
	Timestamp time.Time
}

type StartEvent struct {
	EventID   int
	StartTime time.Time
	Name      string
}

type SuspendEvent struct {
	EventID     int
	SuspendTime time.Time
}

type MarketResult struct {
	MarketID        int
	WinningOutcomes int
	LosingOutcomes  []int
}

type CloseEvent struct {
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

type CreateMarket struct {
	EventID   int
	MarketID  int
	Name      string
	Outcomes  []CreateMarketOutcome
	Timestamp time.Time
}

type UpdatePrice struct {
	EventID   int
	MarketID  int
	OutcomeID int
	Price     float64
	Timestamp time.Time
}

func (e CreateEvent) GetID() int {
	return e.EventID
}

func (e StartEvent) GetID() int {
	return e.EventID
}

func (e SuspendEvent) GetID() int {
	return e.EventID
}

func (e CloseEvent) GetID() int {
	return e.EventID
}

func (e CreateMarket) GetID() int {
	return e.EventID
}

func (e UpdatePrice) GetID() int {
	return e.EventID
}
