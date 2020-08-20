package main

import (
	"es/projection"
	"es/shell"
	"fmt"

	"github.com/abiosoft/ishell"
)

var sh *ishell.Shell

func RegisterCmd(name string, help string, handler func(*ishell.Context)) {
	sh.AddCmd(&ishell.Cmd{
		Name: name,
		Help: help,
		Func: handler,
	})
}

// json schema for validating commands
func main() {
	sh = ishell.New()

	sh.Println(shell.Green("Event Sourcing shell"))

	RegisterCmd("create_event", "<event id> <name> <sport>", shell.CreateEvent)
	RegisterCmd("start_event", "<event id>", shell.StartEvent)
	RegisterCmd("suspend_event", "<event id>", shell.SuspendEvent)
	RegisterCmd("close_event", "<event id>", shell.CloseEvent)
	RegisterCmd("print_active_events", "List active events", PrintActiveEvents)

	RegisterCmd(
		"create_market",
		"<event id> <market id> <name> [<outcome id> <outcome name> <starting price>]...",
		shell.CreateMarket,
	)

	RegisterCmd("update_price", "<event id> <market id> <outcome id> <price>", shell.UpdatePrice)
	RegisterCmd("print_prices", "Print prices of active events", PrintPrices)

	sh.Process("create_event", "1", "Man-Ars", "soccer")
	sh.Process("create_market", "1", "1", "Win-Draw-Win",
		"1", "Home team", "1.4",
		"2", "Draw", "2.8",
		"3", "Away team", "3.7")

	sh.Run()
}

func PrintPrices(c *ishell.Context) {
	c.Println(shell.Yellow(fmt.Sprintf("%5s %5s %5s %10s %20s", "Event", "Mkt", "Outcome", "Price", "Valid from")))

	for _, p := range projection.ActiveEventPriceView {
		c.Printf("%5d %5d %5d %10.2f %20s\n", p.EventID, p.MarketID, p.OutcomeID, p.Price,
			shell.DateFormat(p.ValidFrom))
	}
}

func PrintActiveEvents(c *ishell.Context) {
	c.Println(shell.Yellow(fmt.Sprintf("%5s %20s %s", "Event", "Name", "Start time")))

	for _, e := range projection.ActiveEvents {
		c.Printf("%5d %20s %s\n", e.ID, e.Name, shell.DateFormat(e.StartTime))
	}
}
