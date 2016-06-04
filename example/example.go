package main

import (
	"fmt"

	wp "github.com/elauqsap/workerpool"
)

// SingleJob is the structure that will be passed for processing
type SingleJob struct {
	id int
}

// Process is an implementation of wp.Job.Process()
func (s SingleJob) Process() {
	fmt.Printf("Job ID: %d\n", s.id)
}

func main() {
	// Initialize a Dispatcher
	dispatcher := wp.NewDispatcher(4, 1024)
	// Start the Dispatcher and create/register the Workers to the WorkerPool
	dispatcher.Run()
	// Queue two jobs for processing
	dispatcher.WaitGroup.Add(4)
	dispatcher.JobQueue <- SingleJob{id: 1}
	dispatcher.JobQueue <- SingleJob{id: 2}
	dispatcher.JobQueue <- SingleJob{id: 3}
	dispatcher.JobQueue <- SingleJob{id: 4}

	// Block main thread until processing in go routines completes
	dispatcher.WaitGroup.Wait()
}
