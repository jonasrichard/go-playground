package main

import (
	"fmt"
	"math/rand"
	"time"
)

const partitions = 4

type WorkItem struct {
	Key   int
	Value string
}

func sleep(base int, variance int) {
	secs := time.Duration(base + rand.Intn(variance))
	time.Sleep(secs * time.Millisecond)
}

func worker(partition int, data <-chan WorkItem) {
    for item := range data {
		fmt.Printf("%v processing %v => %v\n", partition, item.Key, item.Value)

		if partition == 0 {
			sleep(10000, 10000)
		} else {
			sleep(500, 500)
		}

		fmt.Printf("%v finished with %v\n", partition, item.Key)
	}
}

func try_send_item(item WorkItem, workers map[int]chan<- WorkItem) bool {
	partition := item.Key % partitions
	ch, _ := workers[partition]

	select {
	case ch <- item:
		return true
	default:
		return false
	}
}

func try_sending_pending(pending []WorkItem, workers map[int]chan<- WorkItem) []WorkItem {
	var newPending = make([]WorkItem, 0)
	for _, item := range pending {
		if !try_send_item(item, workers) {
			newPending = append(newPending, item)
		} else {
			fmt.Printf("Managed to send to %v the pending item %v\n", item.Key%partitions, item.Key)
		}
	}

	return newPending
}

func router(data <-chan WorkItem) {
	var workers map[int]chan<- WorkItem = make(map[int]chan<- WorkItem)
	var pending []WorkItem = make([]WorkItem, 0)

	for i := 0; i < partitions; i++ {
		ch := make(chan WorkItem)
		workers[i] = ch
		go worker(i, ch)
	}

	for {
		select {
		case item := <-data:
			// this is almost good but we only send data when we got data!
			pending = try_sending_pending(pending, workers)

			if !try_send_item(item, workers) {
				fmt.Printf("%v is blocked adding %v to the pending list\n", item.Key%partitions, item.Key)
				pending = append(pending, item)
			}
		}
	}
}

func main() {
	var data_flow = make(chan WorkItem)

	go router(data_flow)

	for {
		data_flow <- WorkItem{Key: rand.Intn(1024), Value: "aaa"}
		sleep(150, 100)
	}
}
