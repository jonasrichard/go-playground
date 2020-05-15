package main

import (
	"errors"

	"github.com/labstack/gommon/log"
)

type Event struct {
	id int
}

type Market struct {
	id      int
	eventId int
}

type Outcome struct {
	id       int
	eventId  int
	marketId int
}

type Price struct {
	id        int
	eventId   int
	outcomeId int
}

type JobMap map[int][]interface{}

const partitions = 16

var NotFound = errors.New("not found")

func findMapKeyForEvent(active JobMap, eventId int) (int, error) {
	for p := 0; p < partitions; p++ {
		jobs, ok := active[p]
		if !ok {
			continue
		}

		for i := range jobs {
			switch j := jobs[i].(type) {
			case Event:
				if j.id == eventId {
					return p, nil
				}
			case Market:
				if j.eventId == eventId {
					return p, nil
				}
			case Outcome:
				if j.eventId == eventId {
					return p, nil
				}
			case Price:
				if j.eventId == eventId {
					return p, nil
				}
			}
		}
	}

	return 0, NotFound
}

func findDependent(active JobMap, data interface{}) (int, error) {
    switch v := data.(type) {
	case Event:
		log.Infof("Event: %v", v)

        return findMapKeyForEvent(active, v.id)
	}

    return 0, NotFound
}

func appendMap(active JobMap, partition int, data interface{}) {
    arr, ok := active[partition]
    if ok {
        for i := range arr {
            if arr[i] == data {
                return
            }
        }

        arr = append(arr, data)
        active[partition] = arr
    } else {
        arr := make([]interface{}, 1)
        arr[0] = data
        active[partition] = arr
    }
}

func router(input chan interface{}, stop chan bool) {
    active := make(JobMap)

    for data := range input {
        partition, err := findDependent(active, data)

        if err == NotFound {
            log.Infof("No partition found")
        } else {
            log.Infof("Partition is %d", partition)
        }

        appendMap(active, partition, data)
        log.Infof("%v", active)
        log.Infof("%T", data)
    }

    stop <-true
}

func main() {
	var queue = make(chan interface{})
    var ready = make(chan bool)

	go router(queue, ready)

    for i := 0; i < 4; i++ {
        queue <-Event{id: 3}
    }

    close(queue)

    <-ready
}
