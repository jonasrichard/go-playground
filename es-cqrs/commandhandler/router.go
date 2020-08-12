package commandhandler

import (
	"errors"
	"es/command"
	"es/store"
	"fmt"
)

var ErrUnknownCommand = errors.New("unknown command")

// Dynamic registration of handler funcs
func Handle(c interface{}) (interface{}, error) {
	var evt interface{}
	var err error

	switch cmd := c.(type) {
	case command.CreateEventCommand:
		evt, err = HandleCreateEventCommand(cmd)
	case command.StartEventCommand:
		evt, err = HandleStartEventCommand(cmd)
	case command.CreateMarketCommand:
		evt, err = HandleCreateMarketCommand(cmd)
	}

	fmt.Printf("Router %v %v\n", evt, err)

	if err == nil {
		event, ok := evt.(store.SourceableEvent)
		if ok {
			store.Store(event)

			return evt, nil
		}

		return nil, ErrUnknownCommand
	}

	return nil, err
}
