package main

import (
	"es/projection"
	"es/shell"

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

	sh.Println("Event Sourcing shell")

	RegisterCmd("create_event", "<event id> <name> <sport>", shell.CreateEvent)
	RegisterCmd("start_event", "<event id>", shell.StartEvent)

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
	for _, p := range projection.ActiveEventPriceView {
		c.Printf("%5d %5d %5d %10f %20s\n", p.EventID, p.MarketID, p.OutcomeID, p.Price, p.ValidFrom)
	}
}
