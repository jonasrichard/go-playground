package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

type PartitionID int
type EventID int
type MarketID int

type Event struct {
	ID EventID
}

type Market struct {
	ID      MarketID
	EventID EventID
}

type Outcome struct {
	ID       int
	EventID  EventID
	MarketID MarketID
}

type Price struct {
	ID        int
	EventID   EventID
	OutcomeID int
}

type EventAskRequest struct {
	EventID EventID
	Data    chan interface{}
}

type Consumer struct {
	PartitionID      PartitionID
	Queue            []interface{}
	Input            chan interface{}
	Sender           chan interface{}
	Ask              chan EventAskRequest
	ConsumerFinished chan bool
	SenderReady      chan bool
	SenderFinished   chan bool
}

type Partitions map[PartitionID]Consumer

type Output interface {
	Send(data interface{})
}

type NoOp struct {
}

const partitionCount PartitionID = 16

//var NotFound = errors.New("not found")

// TODO probably here we need a Mutex
//var LastIDs map[int]int

//func findMapKeyForEvent(active Partitions, eventId int) (int, error) {
//	for p := 0; p < partitions; p++ {
//		jobs, ok := active[p]
//		if !ok {
//			continue
//		}
//
//		for i := range jobs {
//			switch j := jobs[i].(type) {
//			case Event:
//				if j.id == eventId {
//					return p, nil
//				}
//			case Market:
//				if j.eventId == eventId {
//					return p, nil
//				}
//			case Outcome:
//				if j.eventId == eventId {
//					return p, nil
//				}
//			case Price:
//				if j.eventId == eventId {
//					return p, nil
//				}
//			}
//		}
//	}
//
//	return 0, NotFound
//}
//
//func findDependent(active Partitions, data interface{}) (int, error) {
//	switch v := data.(type) {
//	case Event:
//		log.Infof("Event: %v", v)
//
//		return findMapKeyForEvent(active, v.id)
//	}
//
//	return 0, NotFound
//}
//
//func router(input chan interface{}, stop chan bool) {
//	active := make(Partitions)
//
//    for i := 0; i < partitions; i++ {
//        consumer := Consumer{
//            Queue: make([]interface{}, 0),
//            Input: make(chan interface{}),
//            Ask: make(chan int),
//        }
//
//        active[i] = consumer
//    }
//
//	for data := range input {
//		partition, err := findDependent(active, data)
//
//		if err == NotFound {
//			log.Infof("No partition found")
//		} else {
//			log.Infof("Partition is %d", partition)
//		}
//
//		appendMap(active, partition, data)
//		log.Infof("%v", active)
//		log.Infof("%T", data)
//	}
//
//	stop <- true
//}
//
//func() {
//    for data := stream {
//        switch data.(type) {
//        case Event:
//            if activeMarket[eventId] != partition_id {
//                // we need to steal or wait
//                dirtyEvent[eventId] = true
//            } else {
//                send
//            }
//
//        case Market:
//            if dirtyEvent[market.eventId] || activeMarket[marketId] != partition_id {
//                // wait
//            }
//            activeMarket[marketId] = partition_id
//            activeEvent[market.eventId] = partition_id
//        }
//    }
//}

func (o *NoOp) Send(data interface{}) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func send(input chan interface{}, senderReady chan bool, senderFinished chan bool, output Output) {
	for item := range input {
		log.Infof("Finished %v", item)
		output.Send(item)
		senderReady <- true
	}

	senderFinished <- true
}

var dirtyMu sync.Mutex
var dirtyEvents map[EventID][]PartitionID

func appendDirtyEvent(eventID EventID, partitionID PartitionID) {
	dirtyMu.Lock()
	defer dirtyMu.Unlock()

	partitions, ok := dirtyEvents[eventID]
	if !ok {
		partitions := make([]PartitionID, 1)
		partitions[0] = partitionID
	} else {
		for i := range partitions {
			if partitions[i] == partitionID {
				return
			}
		}

		partitions = append(partitions, partitionID)
	}

	dirtyEvents[eventID] = partitions
}

func removeDirtyEvent(eventID EventID, partitionID PartitionID) {
	dirtyMu.Lock()
	defer dirtyMu.Unlock()

	partitions, ok := dirtyEvents[eventID]
	if !ok {
		return
	}

	for i := range partitions {
		if partitions[i] == partitionID {
			partitions = append(partitions[:i], partitions[i+1:]...)
			dirtyEvents[eventID] = partitions

			return
		}
	}
}

func consume(consumer Consumer) {
loop:
	for {
		select {
		case item, ok := <-consumer.Input:
			if !ok {
				log.Infof("%d emptying queue %v", consumer.PartitionID, consumer.Queue)

				for q := range consumer.Queue {
					consumer.Sender <- q
				}

				break loop
			}

			select {
			case consumer.Sender <- item:
				// managed to send
			default:
				log.Infof("Something queued %v", item)
				consumer.Queue = append(consumer.Queue, item)

				switch i := item.(type) {
				case Event:
					appendDirtyEvent(i.ID, consumer.PartitionID)
				}
			}

		case askRequest := <-consumer.Ask:
			log.Infof("Event id %v", askRequest.EventID)
			// event is asked, send back things on the data channel

		case <-consumer.SenderReady:
			// we can send data, if queue is empty no-op

			if len(consumer.Queue) > 0 {
				item := consumer.Queue[0]
				consumer.Queue = consumer.Queue[1:]

				consumer.Sender <- item

				switch i := item.(type) {
				case Event:
					removeDirtyEvent(i.ID, consumer.PartitionID)
				}
			}
		}
	}

	log.Infof("Closing sender #%d", consumer.PartitionID)

	close(consumer.Sender)

	consumer.ConsumerFinished <- true
}

func router(partitions Partitions, input chan interface{}, ready chan bool, output Output) {
	var i PartitionID

	for i = 0; i < partitionCount; i++ {
		consumer := Consumer{
			PartitionID:      PartitionID(i),
			Queue:            make([]interface{}, 0),
			Input:            make(chan interface{}),
			Sender:           make(chan interface{}),
			Ask:              make(chan EventAskRequest),
			ConsumerFinished: make(chan bool),
			SenderReady:      make(chan bool),
			SenderFinished:   make(chan bool),
		}

		partitions[i] = consumer

		go consume(consumer)
		go send(consumer.Sender, consumer.SenderReady, consumer.SenderFinished, output)
	}

	for item := range input {
		log.Infof("Item arrived %v", item)

		var partitionID PartitionID
		switch i := item.(type) {
		case Event:
			partitionID = PartitionID(i.ID) % partitionCount
		case Market:
			partitionID = PartitionID(i.EventID) % partitionCount
		}

		partitions[partitionID].Input <- item
	}

	log.Info("Closing consumers")

	for i := range partitions {
		close(partitions[i].Input)
	}

	for i := range partitions {
		<-partitions[i].ConsumerFinished
		<-partitions[i].SenderFinished
	}

	ready <- true
}

func init() {
	log.SetHeader("${time_rfc3339_nano} ${level} ${short_file}:${line} ${prefix}")
}

func main() {
	//	var queue = make(chan interface{})
	//	var ready = make(chan bool)
	//
	//	go router(queue, ready)
	//
	//	close(queue)
	//
	//	<-ready
}
