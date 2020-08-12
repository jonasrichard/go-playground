package commandhandler_test

import (
	"es/aggregate"
	"es/command"
	"es/commandhandler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventDoubleStart(t *testing.T) {
	cmds := []interface{}{
		command.CreateEventCommand{EventID: 8, Name: "AC Milan - Bari", Type: "Soccer"},
		command.StartEventCommand{EventID: 8},
		command.StartEventCommand{EventID: 8},
	}

	for i, cmd := range cmds {
		_, err := commandhandler.Handle(cmd)

		if i == 2 {
			assert.NotNil(t, err)
		}
	}

	event, _ := aggregate.LoadEvent(8)

	assert.NotNil(t, event)
	assert.Equal(t, event.ID, 8)
}
