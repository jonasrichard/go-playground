package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	// The input of the processing
	inputID int
	// The result of the processing
	result string
	// Storing possible error
	err error
	// Number of retries so far
	retries int
}

const numberOfWorkers = 8

var processingError error = errors.New("some error")

func processJob(job Job, feedback chan Job) {
	time.Sleep(time.Second + time.Duration(rand.Intn(500))*time.Millisecond)

	if 400 > rand.Intn(500) {
		job.err = processingError
	} else {
		job.result = "ready"
	}

	feedback <- job
}

func workerLoop(workerID int, allJobs chan Job, feedback chan Job) {
	fmt.Printf("[worker #%d] Started\n", workerID)

	// This loop finishes when allJobs (so jobStream) is closed
	for job := range allJobs {
		fmt.Printf("[worker #%d] Processing job: %d\n", workerID, job.inputID)

		processJob(job, feedback)
	}

	fmt.Printf("[worker #%d] Finished\n", workerID)
}

func randomIDs(n int) []int {
	ids := make([]int, n)

	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(100)
	}

	return ids
}

// inputProcessor sending the jobs into jobStream if possible. It can do that if workers are
// selecting. If not we have a chance to process jobs to be retried.
// When the processing finished, we can close feedback stream to tell the other goroutine
// that it can stop, too.
func inputProcessor(jobs []Job, jobStream chan Job, retryStream chan Job, feedback chan Job) {
	var cursor int

	for {
		if cursor >= len(jobs) {
			break
		}

		job := jobs[cursor]

		select {
		case jobStream <- job:
			fmt.Printf("[input processor] Job is being sent for processing: %d\n", job.inputID)
			cursor++
		case toRetry := <-retryStream:
			fmt.Printf("[input processor] Job is requeued: %d\n", job.inputID)

			toRetry.retries++
			toRetry.err = nil

			jobs = append(jobs, toRetry)
		}
	}

	// Let workers terminate
	close(jobStream)

	// The only reason we have the feedback here is to close it
	close(feedback)
}

func main() {
	inputs := randomIDs(32)
	// Channel for sending jobs to workers
	jobStream := make(chan Job)
	// Via feedback workers can tell the result of the compuation (or error)
	feedback := make(chan Job)
	// Channel for feeding back jobs to retry
	retryStream := make(chan Job)

	for i := 0; i < numberOfWorkers; i++ {
		go workerLoop(i, jobStream, feedback)
	}

	jobs := make([]Job, len(inputs))
	for i := 0; i < len(jobs); i++ {
		jobs[i] = Job{inputID: inputs[i]}
	}

	go inputProcessor(jobs, jobStream, retryStream, feedback)

	for result := range feedback {
		fmt.Printf("[main] Feedback: %s (%d)\n", result.err, result.inputID)

		if result.retries < 5 {
			retryStream <- result
		}
	}
}
