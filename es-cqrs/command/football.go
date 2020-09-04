package command

import "time"

type CreateFootballEventCommand struct {
	EventID            int
	HomeTeam           string
	AwayTeam           string
	StartTime          time.Time
	FirstGoalMarket    bool
	HowManyGoalsMarket int // -1 no market generated
}
