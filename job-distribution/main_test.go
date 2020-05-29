package main

import (
	"sync"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

type Collector struct {
	mu   sync.Mutex
	data []interface{}
}

func (c *Collector) Send(item interface{}) {
    c.mu.Lock()
	c.data = append(c.data, item)
    c.mu.Unlock()
}

func TestEvent(t *testing.T) {
	var partitions Partitions = make(Partitions)
	var queue = make(chan interface{})
	var ready = make(chan bool)
	var collector Collector

	go router(partitions, queue, ready, &collector)

	for i := 0; i < 32; i++ {
		queue <- Event{ID: EventID(i)}
	}

	close(queue)

	<-ready

	log.Infof("%v", collector.data)

	assert.Equal(t, len(collector.data), 32)
}
