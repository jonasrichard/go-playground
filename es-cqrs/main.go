package main

import (
	"es/command"
	"es/commandhandler"
	"es/projection"
	"es/store"
	"fmt"
	"strconv"

	"github.com/abiosoft/ishell"
)

func SimpleHandler(command interface{}) error {
	fmt.Printf("[repl] %v\n", command)

	event, err := commandhandler.Handle(command)

	if err == nil {
		fmt.Printf("[repl] Event: %v\n", event)

		return commandhandler.EventBus(event.(store.SourceableEvent))
	}

	fmt.Printf("[repl] Error %v\n", err)

	return err
}

// json schema for validating commands
func main() {
	shell := ishell.New()

	shell.Println("Event Sourcing shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "create_event",
		Help: "Create event",
		Func: func(c *ishell.Context) {
			eventID, err := strconv.Atoi(c.Args[0])
			if err != nil {
				c.Println(err)

				return
			}

			cmd := command.CreateEventCommand{EventID: eventID, Name: c.Args[1], Type: c.Args[2]}

			if err := SimpleHandler(cmd); err != nil {
				c.Println(err)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "create_market",
		Help: "Create market with outcomes",
		Func: func(c *ishell.Context) {
			eventID, _ := strconv.Atoi(c.Args[0])
			marketID, _ := strconv.Atoi(c.Args[1])

			outcomes := []command.CreateMarketOutcome{}
			i := 3
			for {
				if i >= len(c.Args) {
					break
				}

				outcomeID, _ := strconv.Atoi(c.Args[i])
				price, _ := strconv.ParseFloat(c.Args[i+2], 64)

				outcome := command.CreateMarketOutcome{ID: outcomeID, Name: c.Args[i+1], StartingPrice: price}
				outcomes = append(outcomes, outcome)

				i += 3
			}

			cmd := command.CreateMarketCommand{
				EventID:  eventID,
				MarketID: marketID,
				Name:     c.Args[2],
				Outcomes: outcomes,
			}

			if err := SimpleHandler(cmd); err != nil {
				c.Println(err)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "print_prices",
		Help: "Print prices of active events",
		Func: func(c *ishell.Context) {
			for _, p := range projection.ActiveEventPriceView {
				c.Printf("%5d %5d %5d %10f %20s\n", p.EventID, p.MarketID, p.OutcomeID, p.Price, p.ValidFrom)
			}
		},
	})

	shell.Run()
}
