package event

import "time"

type CreateEvent struct {
	EventID   int
	Name      string
	Type      string
	Timestamp time.Time
}

func (e CreateEvent) GetID() int {
	return e.EventID
}

type StartEvent struct {
	EventID   int
	StartTime time.Time
}

func (e StartEvent) GetID() int {
	return e.EventID
}

type SuspendEvent struct {
	EventID     int
	SuspendTime time.Time
}

func (e SuspendEvent) GetID() int {
	return e.EventID
}

type CloseEvent struct {
	EventID   int
	CloseTime time.Time
	// Resulting?
}

func (e CloseEvent) GetID() int {
	return e.EventID
}

type CreateMarket struct {
	EventID  int
	MarketID int
	Name     string
	Outcomes []struct {
		ID            int
		Name          string
		StartingPrice float32
	}
	Timestamp time.Time
}

func (e CreateMarket) GetID() int {
	return e.EventID
}

type UpdatePrice struct {
	EventID   int
	MarketID  int
	OutcomeID int
	Price     float32
	Timestamp time.Time
}

func (e UpdatePrice) GetID() int {
	return e.EventID
}