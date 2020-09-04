package commandhandler

import (
	"es/aggregate"
	"es/command"
	"es/event"
	"es/store"
	"fmt"
	"time"
)

func HandleCreateFootballEventCommand(cmd command.CreateFootballEventCommand) ([]store.SourceableEvent, error) {
	_, err := aggregate.LoadEvent(cmd.EventID)

	if err == nil {
		return nil, ErrEventAlreadyCreated
	}

	created := event.EventCreated{
		EventID:   cmd.EventID,
		Name:      cmd.HomeTeam + "-" + cmd.AwayTeam,
		Type:      "soccer",
		Timestamp: time.Now(),
	}

	wdw := event.MarketCreated{
		EventID:  cmd.EventID,
		MarketID: 0, // generate an id
		Name:     "Win/Draw/Win",
		Outcomes: []event.CreateMarketOutcome{
			{0, cmd.HomeTeam + " wins", 0.0},
			{0, "Draw", 0.0},
			{0, cmd.AwayTeam + " wins", 0.0},
		},
		Timestamp: time.Now(),
	}

	result := []store.SourceableEvent{created, wdw}

	if cmd.FirstGoalMarket {
		firstGoal := event.MarketCreated{
			EventID:  cmd.EventID,
			MarketID: 0,
			Name:     "First goal scored",
			Outcomes: []event.CreateMarketOutcome{
				{0, cmd.HomeTeam + " scores first", 0.0},
				{0, cmd.AwayTeam + " scores first", 0.0},
			},
			Timestamp: time.Now(),
		}

		result = append(result, firstGoal)
	}

	if cmd.HowManyGoalsMarket > -1 {
		outcomes := make([]event.CreateMarketOutcome, cmd.HowManyGoalsMarket+1)

		for i := 0; i <= cmd.HowManyGoalsMarket; i++ {
			outcome := event.CreateMarketOutcome{
				ID:            0,
				Name:          fmt.Sprintf("%d goal(s) scored", i),
				StartingPrice: 0.0,
			}

			outcomes = append(outcomes, outcome)
		}

		goals := event.MarketCreated{
			EventID:   cmd.EventID,
			MarketID:  0,
			Name:      "How many goals scored",
			Outcomes:  outcomes,
			Timestamp: time.Now(),
		}

		result = append(result, goals)
	}

	return result, nil
}
