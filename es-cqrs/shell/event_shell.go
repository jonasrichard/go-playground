package shell

import (
	"es/command"
	"es/commandhandler"
	"es/store"
	"fmt"
	"strconv"

	"github.com/abiosoft/ishell"
)

func SimpleHandler(command interface{}) {
	fmt.Printf("[repl] %v\n", command)

	event, err := commandhandler.Handle(command)

	if err == nil {
		fmt.Printf("[repl] Event: %v\n", event)

		commandhandler.EventBus(event.(store.SourceableEvent))
	}

	fmt.Printf("[repl] Error %v\n", err)
}

func ToID(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return id
}

func CreateEvent(c *ishell.Context) {
	eventID := ToID(c.Args[0])
	cmd := command.CreateEventCommand{EventID: eventID, Name: c.Args[1], Type: c.Args[2]}

	SimpleHandler(cmd)
}

func StartEvent(c *ishell.Context) {
	eventID := ToID(c.Args[0])
	cmd := command.StartEventCommand{EventID: eventID}

	SimpleHandler(cmd)
}
