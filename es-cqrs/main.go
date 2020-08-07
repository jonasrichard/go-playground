package main

import (
	"es/command"
	"es/commandhandler"
	"fmt"
)

func main() {
	cmd := command.CreateEventCommand{ID: 2, Name: "Man-Ars", Type: "Soccer"}

	event, _ := commandhandler.HandleCreateEventCommand(cmd)

	fmt.Printf("%v\n", event)
}
