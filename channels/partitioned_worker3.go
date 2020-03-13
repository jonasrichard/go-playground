package main

import (
	"fmt"
	"math/rand"
	"reflect"
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

func worker(partition int, data <-chan WorkItem, control chan bool) {
	for {
		select {
		case item := <-data:
			fmt.Printf("%v processing %v => %v\n", partition, item.Key, item.Value)
			if partition == 0 {
				sleep(10000, 10000)
			} else {
				sleep(500, 500)
			}
			fmt.Printf("%v finished with %v\n", partition, item.Key)
			control <- true
		}
	}
}

func find_item_for_partition(items []WorkItem, partition int, pendingItem *WorkItem) ([]WorkItem, bool) {
	fmt.Printf("Old pending: %v\n", items)
	for i, item := range items {
		if item.Key%partitions == partition {
			newPending := append(items[:i], items[i+1:]...)
			fmt.Printf("New pending: %v\n", newPending)
			*pendingItem = item
			return newPending, true
		}
	}

	return items, false
}

func router(data <-chan WorkItem) {
	var workers []chan<- WorkItem = make([]chan<- WorkItem, partitions)
	var control []chan bool = make([]chan bool, partitions)
	var pending []WorkItem = make([]WorkItem, 0)
	var available []bool = make([]bool, partitions)

	for i := 0; i < partitions; i++ {
		ch := make(chan WorkItem)
		ctrl := make(chan bool)
		workers[i] = ch
		control[i] = ctrl
		go worker(i, ch, ctrl)
		available[i] = true
	}

	// For the multiselect we construct workers(i) + data channel
	channels := make([]reflect.SelectCase, partitions+1)
	for i, ctrl := range control {
		channels[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctrl)}
	}
	channels[partitions] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(data)}

	for {
		chosen, value, ok := reflect.Select(channels)

		if !ok {
			return
		}

		if chosen == partitions {
			// we get new work item
			item := value.Interface().(WorkItem)
			partition := item.Key % partitions

			if available[partition] {
				workers[partition] <- item
				available[partition] = false
			} else {
				fmt.Printf("%v is blocked, %v added to the pending list\n", partition, item.Key)
				pending = append(pending, item)
				fmt.Println(pending)
			}
		} else {
			// a worker has become available
			fmt.Printf("%v became available\n", chosen)
			fmt.Println(pending)

			var pendingItem WorkItem
			pending, ok = find_item_for_partition(pending, chosen, &pendingItem)

			if ok {
				fmt.Printf("%v is processing %v pending item\n", chosen, pendingItem.Key)
				workers[chosen] <- pendingItem
			} else {
				available[chosen] = true
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
