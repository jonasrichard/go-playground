package shell

import (
	"es/command"
	"strconv"

	"github.com/abiosoft/ishell"
)

func CreateMarket(c *ishell.Context) {
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

	SimpleHandler(cmd)
}
