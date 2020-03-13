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
		}
	}
}

func router(data <-chan WorkItem) {
	var workers map[int]chan<- WorkItem = make(map[int]chan<- WorkItem)

	for i := 0; i < partitions; i++ {
		ch := make(chan WorkItem)
		workers[i] = ch
		go worker(i, ch)
	}

	for {
		select {
		case item := <-data:
			partition := item.Key % partitions

			ch, _ := workers[partition]
			// TODO this is blocking if the worker is busy!
			// even if we have available workers :(
			ch <- item
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
