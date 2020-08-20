package commandhandler

import (
	"errors"
	"es/aggregate"
	"es/command"
	"es/event"
	"time"
)

var ErrEventAlreadyCreated = errors.New("event exists with that ID")
var ErrEventAlreadyStarted = errors.New("event cannot be started")
var ErrEventAlreadySuspended = errors.New("event cannot be suspended")
var ErrEventAlreadyClosed = errors.New("event cannot be closed")
var ErrEventNotFound = errors.New("event not found")
var ErrCannotChangePrice = errors.New("cannot change price")

func HandleCreateEventCommand(cmd command.CreateEventCommand) (event.CreateEvent, error) {
	_, err := aggregate.LoadEvent(cmd.EventID)

	if err == nil {
		return event.CreateEvent{}, ErrEventAlreadyCreated
	}

	createEvent := event.CreateEvent{
		EventID:   cmd.EventID,
		Name:      cmd.Name,
		Type:      cmd.Type,
		Timestamp: time.Now(),
	}

	return createEvent, nil
}

func HandleStartEventCommand(cmd command.StartEventCommand) (event.StartEvent, error) {
	var result event.StartEvent

	evt, err := aggregate.LoadEvent(cmd.EventID)

	if err != nil {
		return result, err
	}

	if evt.State == aggregate.PreGame ||
		evt.State == aggregate.Suspended {
		startEvent := event.StartEvent{
			EventID:   cmd.EventID,
			StartTime: time.Now(),
			Name:      evt.Name,
		}

		return startEvent, nil
	}

	return result, ErrEventAlreadyStarted
}

func HandleSuspendEventCommand(cmd command.SuspendEventCommand) (event.SuspendEvent, error) {
	eventAggregate, _ := aggregate.LoadEvent(cmd.EventID)

	if eventAggregate.State == aggregate.Started {
		suspendEvent := event.SuspendEvent{
			EventID:     cmd.EventID,
			SuspendTime: time.Now(),
		}

		return suspendEvent, nil
	}

	return event.SuspendEvent{}, ErrEventAlreadySuspended
}

func HandleCloseEventCommand(cmd command.CloseEventCommand) (event.CloseEvent, error) {
	eventAggregate, _ := aggregate.LoadEvent(cmd.EventID)

	if eventAggregate.State != aggregate.Closed {
		closeEvent := event.CloseEvent{
			EventID:   cmd.EventID,
			CloseTime: time.Now(),
		}

		return closeEvent, nil
	}

	return event.CloseEvent{}, ErrEventAlreadyClosed
}

func HandleCreateMarketCommand(cmd command.CreateMarketCommand) (event.CreateMarket, error) {
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

		createMarket := event.CreateMarket{
			EventID:   cmd.EventID,
			MarketID:  cmd.MarketID,
			Name:      cmd.Name,
			Outcomes:  outcomes,
			Timestamp: time.Now(),
		}

		return createMarket, nil
	}

	return event.CreateMarket{}, ErrEventNotFound
}

func HandleUpdatePriceCommand(cmd command.UpdatePriceCommand) (event.UpdatePrice, error) {
	evt, _ := aggregate.LoadEvent(cmd.EventID)

	if evt.State != aggregate.Closed {
		updatePrice := event.UpdatePrice{
			EventID:   cmd.EventID,
			MarketID:  cmd.MarketID,
			OutcomeID: cmd.OutcomeID,
			Price:     cmd.Price,
			Timestamp: time.Now(),
		}

		return updatePrice, nil
	}

	return event.UpdatePrice{}, ErrCannotChangePrice
}
