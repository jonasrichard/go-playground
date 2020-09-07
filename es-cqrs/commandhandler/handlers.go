package commandhandler

import (
	"errors"
	"es/aggregate"
	"es/command"
	"es/event"
	"es/store"
	"time"
)

var ErrEventAlreadyCreated = errors.New("event exists with that ID")
var ErrEventAlreadyStarted = errors.New("event cannot be started")
var ErrEventAlreadySuspended = errors.New("event cannot be suspended")
var ErrEventAlreadyClosed = errors.New("event cannot be closed")
var ErrEventNotFound = errors.New("event not found")
var ErrCannotChangePrice = errors.New("cannot change price")

func HandleCreateEventCommand(cmd command.CreateEventCommand) ([]store.SourceableEvent, error) {
	_, err := aggregate.LoadEvent(cmd.EventID)

	if err == nil {
		return nil, ErrEventAlreadyCreated
	}

	createEvent := event.EventCreated{
		EventID:          cmd.EventID,
		Name:             cmd.Name,
		Type:             cmd.Type,
		PlannedStartTime: time.Now().Add(15 * time.Second),
		Timestamp:        time.Now(),
	}

	return []store.SourceableEvent{createEvent}, nil
}

func HandleUpdateEventCommand(cmd command.UpdateEventCommand) ([]store.SourceableEvent, error) {
	_, err := aggregate.LoadEvent(cmd.EventID)

	if err == nil {
		return nil, ErrEventNotFound
	}

	updateEvent := event.EventUpdated{
		EventID: cmd.EventID,
		Name:    cmd.Name,
		Type:    cmd.Type,
	}

	return []store.SourceableEvent{updateEvent}, nil
}

func HandleStartEventCommand(cmd command.StartEventCommand) ([]store.SourceableEvent, error) {
	evt, err := aggregate.LoadEvent(cmd.EventID)

	if err != nil {
		return nil, err
	}

	if evt.State == aggregate.PreGame ||
		evt.State == aggregate.Suspended {
		startEvent := event.EventStarted{
			EventID:   cmd.EventID,
			StartTime: time.Now(),
			Name:      evt.Name,
		}

		return []store.SourceableEvent{startEvent}, nil
	}

	return nil, ErrEventAlreadyStarted
}

func HandleSuspendEventCommand(cmd command.SuspendEventCommand) ([]store.SourceableEvent, error) {
	eventAggregate, _ := aggregate.LoadEvent(cmd.EventID)

	if eventAggregate.State == aggregate.Started {
		suspendEvent := event.EventSuspended{
			EventID:     cmd.EventID,
			SuspendTime: time.Now(),
		}

		return []store.SourceableEvent{suspendEvent}, nil
	}

	return nil, ErrEventAlreadySuspended
}

func HandleCloseEventCommand(cmd command.CloseEventCommand) ([]store.SourceableEvent, error) {
	evt, _ := aggregate.LoadEvent(cmd.EventID)

	if evt.State != aggregate.Closed {
		// TODO update the markets and outcome which are winning and losing
		closeEvent := event.EventClosed{
			EventID:   cmd.EventID,
			CloseTime: time.Now(),
		}

		return []store.SourceableEvent{closeEvent}, nil
	}

	return nil, ErrEventAlreadyClosed
}

func HandleCreateMarketCommand(cmd command.CreateMarketCommand) ([]store.SourceableEvent, error) {
	_, err := aggregate.LoadEvent(cmd.EventID)

	if err == nil {
		outcomes := []event.CreateMarketOutcome{}

		for _, o := range cmd.Outcomes {
			outcome := event.CreateMarketOutcome{
				ID:            o.ID,
				Name:          o.Name,
				StartingPrice: o.StartingPrice,
			}

			outcomes = append(outcomes, outcome)
		}

		createMarket := event.MarketCreated{
			EventID:   cmd.EventID,
			MarketID:  cmd.MarketID,
			Name:      cmd.Name,
			Outcomes:  outcomes,
			Timestamp: time.Now(),
		}

		return []store.SourceableEvent{createMarket}, nil
	}

	return nil, ErrEventNotFound
}

func HandleUpdatePriceCommand(cmd command.UpdatePriceCommand) ([]store.SourceableEvent, error) {
	evt, _ := aggregate.LoadEvent(cmd.EventID)

	if evt.State != aggregate.Closed {
		updatePrice := event.PriceUpdated{
			EventID:   cmd.EventID,
			MarketID:  cmd.MarketID,
			OutcomeID: cmd.OutcomeID,
			Price:     cmd.Price,
			Timestamp: time.Now(),
		}

		return []store.SourceableEvent{updatePrice}, nil
	}

	return nil, ErrCannotChangePrice
}
