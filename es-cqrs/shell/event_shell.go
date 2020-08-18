package shell

import (
	"es/command"
	"es/commandhandler"
	"es/store"
	"fmt"
	"strconv"

	"github.com/abiosoft/ishell"
	"github.com/lucasepe/color"
)

var Green = color.New(color.FgGreen).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

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

func SuspendEvent(c *ishell.Context) {
	eventID := ToID(c.Args[0])
	cmd := command.SuspendEventCommand{EventID: eventID}

	SimpleHandler(cmd)
}

func CloseEvent(c *ishell.Context) {
	eventID := ToID(c.Args[0])
	cmd := command.CloseEventCommand{EventID: eventID}

	SimpleHandler(cmd)
}
