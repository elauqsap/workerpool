package workerpool

import "sync"

// Job interface must be implemented by each "job type" in order to be
// processed by the worker pool
type Job interface {
	Process()
}

// Dispatcher creates workers and dispatches jobs when received
type Dispatcher struct {
	JobQueue   chan Job
	MaxWorkers int
	WaitGroup  *sync.WaitGroup
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
}

// NewDispatcher creates a dispatcher that is used to create workers
// and dispatch jobs to them
func NewDispatcher(maxWorkers int, queueSize int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{JobQueue: make(chan Job, queueSize), MaxWorkers: maxWorkers, WorkerPool: pool, WaitGroup: &sync.WaitGroup{}}
}

// Run creates the workers and dispatches jobs from a JobQueue to each worker
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.WaitGroup)
		worker.Start()
	}

	// start the dispatcher routine
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

// Worker represents the worker that executes the job
type Worker struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	// A channel for receiving a job that was dispatched
	JobChannel chan Job
	// A channel for receiving a worker termination signal
	// (quits after processing)
	quit chan bool
	// A WaitGroup to signal the completed processing of a Job
	wg *sync.WaitGroup
}

// NewWorker creates a new worker that can be registered to a WorkerPool
// and receive jobs
func NewWorker(workerPool chan chan Job, wg *sync.WaitGroup) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		wg:         wg}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Process()
				// signal to the wait group that a queued job has been processed
				// so the main thread can continue
				w.wg.Done()
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}
