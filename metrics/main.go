package main

import (
	//	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"runtime/pprof"

	"github.com/rcrowley/go-metrics"
	"github.com/vrischmann/go-metrics-influxdb"
)

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go processOne()

	go metrics.Log(metrics.DefaultRegistry, 2*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	go influxdb.InfluxDBWithTags(
		metrics.DefaultRegistry,
		2*time.Second,
		"http://127.0.0.1:8086",
		"go_metrics",
		"app.counter",
		"admin",
		"admin",
		make(map[string]string, 0),
		true)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
}

func processOne() {
	counter := metrics.NewCounter()
	metrics.Register("app.counter", counter)

	sample := metrics.NewExpDecaySample(1028, 0.015)
	histogram := metrics.NewHistogram(sample)
	metrics.Register("app.histogram", histogram)

	for {
		sleepTime := rand.Intn(500)
		counter.Inc(int64(10 - rand.Intn(18)))
		histogram.Update(int64(sleepTime))

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}
