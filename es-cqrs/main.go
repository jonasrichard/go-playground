package main

import (
	"es/command"
	"es/commandhandler"
	"fmt"
)

// json schema for validating commands
func main() {
	cmd := command.CreateEventCommand{EventID: 2, Name: "Man-Ars", Type: "Soccer"}

	event, _ := commandhandler.HandleCreateEventCommand(cmd)

	fmt.Printf("%v\n", event)
}
