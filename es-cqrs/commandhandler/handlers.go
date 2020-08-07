package commandhandler

import (
	"es/command"
	"es/event"
	"es/eventstore"
)

func HandleCreateEventCommand(cmd command.CreateEventCommand) (event.CreateEventEvent, error) {
	event := event.CreateEventEvent{
		AggregateID: cmd.ID,
		EventID:     1,
		Name:        cmd.Name,
		Type:        cmd.Type,
	}

	eventstore.Store(event)

	return event, nil
}
