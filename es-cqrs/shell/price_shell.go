package shell

import (
	"es/command"
	"strconv"

	"github.com/abiosoft/ishell"
)

func UpdatePrice(c *ishell.Context) {
	eventID := ToID(c.Args[0])
	marketID := ToID(c.Args[1])
	outcomeID := ToID(c.Args[2])
	price, _ := strconv.ParseFloat(c.Args[3], 64)

	cmd := command.UpdatePriceCommand{
		EventID:   eventID,
		MarketID:  marketID,
		OutcomeID: outcomeID,
		Price:     price,
	}

	SimpleHandler(cmd)
}
