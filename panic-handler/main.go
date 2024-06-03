package main

import (
	"log"
	"time"
)

func Producer(exited chan<- bool) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Producer recovered from panic: %v\n", r)
        }

        close(exited)
    }()

    i := 0

	for {
		time.Sleep(250 * time.Millisecond)

		log.Printf("Producer\n")

        i += 1
        if i > 3 {
            // Normally it would kill main goroutine
            panic("Producer paniced before normal return")
        }

        if i > 10 {
            exited <- true

            return
        }
	}
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Main recovered from panic: %v\n", r)
        }
    }()

	exited := make(chan bool)

	go Producer(exited)

    for range exited {
    }

    log.Printf("Main completed normally")
}
