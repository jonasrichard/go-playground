package commandhandler

import (
	"errors"
	"es/command"
	"es/projection"
	"es/store"
)

var ErrUnknownCommand = errors.New("unknown command")

func Handle(cmd interface{}) ([]store.SourceableEvent, error) {
	events, err := ExecuteCommand(cmd)

	if err == nil {
		for _, evt := range events {
			store.Store(evt)
		}
	}

	return events, err
}

func ExecuteCommand(c interface{}) ([]store.SourceableEvent, error) {
	switch cmd := c.(type) {
	case command.CreateEventCommand:
		return HandleCreateEventCommand(cmd)
	case command.StartEventCommand:
		return HandleStartEventCommand(cmd)
	case command.SuspendEventCommand:
		return HandleSuspendEventCommand(cmd)
	case command.CloseEventCommand:
		return HandleCloseEventCommand(cmd)
	case command.CreateMarketCommand:
		return HandleCreateMarketCommand(cmd)
	case command.UpdatePriceCommand:
		return HandleUpdatePriceCommand(cmd)
	case command.CreateFootballEventCommand:
		return HandleCreateFootballEventCommand(cmd)
	}

	return nil, ErrUnknownCommand
}

func EventBus(event store.SourceableEvent) error {
	// TODO loop
	projection.ProjectActiveEventPrice(event)
	projection.ProjectActiveEvents(event)
	projection.ProjectMarketCountryView(event)

	return nil
}
