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

	eventAggregate, err := aggregate.LoadEvent(cmd.EventID)

	if err != nil {
		return result, err
	}

	if eventAggregate.State == aggregate.PreGame {
		startEvent := event.StartEvent{
			EventID:   cmd.EventID,
			StartTime: time.Now(),
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
	// find event by id
	// if event is closed we cannot create markets

	return event.CreateMarket{}, nil
}

func HandleUpdatePriceCommand(cmd command.UpdatePriceCommand) (event.UpdatePrice, error) {
	// look for the event by outcome id
	// if event is close we cannot modify prices
	// NOTE: we need to do it in a projection

	return event.UpdatePrice{}, nil
}