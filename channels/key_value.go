package main

// A simple key-value store as goroutine

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	key   int
	value string
}

type ReadOp struct {
	key    int
	result chan Item
}

type WriteOp struct {
	item   Item
	result chan bool
}

func StartKeyValueStore(input <-chan interface{}) {
	var store map[int]string = make(map[int]string, 0)

	for op := range input {
		switch op.(type) {
		case ReadOp:
			readOp := op.(ReadOp)
			fmt.Printf("Reading the value %v\n", readOp.key)

			value, ok := store[readOp.key]
			if ok {
				readOp.result <- Item{readOp.key, value}
			} else {
				readOp.result <- Item{0, ""}
			}
		case WriteOp:
			writeOp := op.(WriteOp)
			fmt.Printf("Writing the item %v\n", writeOp.item)

			store[writeOp.item.key] = writeOp.item.value

			writeOp.result <- true
		default:
			fmt.Printf("Unknown operation %v\n", op)
		}
	}
}

func main() {
	var input = make(chan interface{})

	var readResult = make(chan Item)
	var writeResult = make(chan bool)

	go StartKeyValueStore(input)

	for {
		key := rand.Intn(32)
		switch rand.Intn(2) {
		case 0:
			input <- ReadOp{key: key, result: readResult}
			result := <-readResult
			fmt.Printf("Reading %v... result is %v\n", key, result)
		case 1:
			input <- WriteOp{item: Item{key: key, value: "Apple"}, result: writeResult}
			result := <-writeResult
			fmt.Printf("Writing %v... result is %v\n", key, result)
		}
		time.Sleep(250 * time.Millisecond)
	}
}
