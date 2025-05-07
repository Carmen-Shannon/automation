package worker

import (
	"sync"
	"time"
)

// NewWorker creates a new worker with the given ID, task channel, stop channel, idle timeout, and exit callback function.
// This will return an interface which can be used to manipulate the worker.
//
// Parameters:
//   - id: The ID of the worker. This is an integer value that uniquely identifies the worker.
//   - taskChan: The channel for tasks to be processed by the worker. This is a channel of Task type that the worker will listen to for incoming tasks.
//   - stopChan: The channel for stopping the worker. This is a channel of int type that the worker will listen to for stop signals.
//   - idleTimeout: The timeout duration for the worker to wait before stopping. This is a time.Duration value that determines how long the worker should wait before stopping if there are no tasks.
//   - onExit: The callback function to be called when the worker exits. This is a function that takes an integer parameter (the worker ID) and returns nothing.
func NewWorker(id int, taskChan chan Task, stopChan chan int, idleTimeout time.Duration, onExit func(int)) Worker {
	return &worker{
		mu:          sync.Mutex{},
		id:          id,
		taskChan:    taskChan,
		stopChan:    stopChan,
		idleTimeout: idleTimeout,
		onExit:      onExit,
		active:      false,
	}
}

type worker struct {
	mu sync.Mutex

	id     int
	active bool

	idleTimeout time.Duration
	onExit      func(int)

	taskChan chan Task
	stopChan chan int
}

// Worker is the interface that defines the methods for a worker in the worker pool.
type Worker interface {
	// ID returns the ID of the worker.
	//
	// Returns:
	//   - int: The ID of the worker. This is an integer value that uniquely identifies the worker.
	ID() int

	// IsActive returns true if the worker is active, false otherwise.
	//
	// Returns:
	//   - bool: True if the worker is active, false otherwise.
	IsActive() bool

	// Start starts the worker and begins processing tasks from the task channel.
	// The worker controls it's own lifecycle and will stop when it finished processing tasks and it idles for long enough to reach it's idle timeout threshold.
	Start()

	// Stop stops the worker and cleans up any resources used by the worker.
	Stop()
}

func (w *worker) ID() int {
	return w.id
}

func (w *worker) IsActive() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.active
}

func (w *worker) Start() {
	w.mu.Lock()
	w.active = true
	w.mu.Unlock()
	go func() {
		for {
			select {
			case i, ok := <-w.stopChan:
				if !ok {
					return
				}
				if i == w.id {
					return
				}
			case t, ok := <-w.taskChan:
				if !ok {
					return
				}

				t.Do()
			}
		}
	}()
}

func (w *worker) Stop() {
	w.stopChan <- w.id
	w.mu.Lock()
	w.active = false
	w.mu.Unlock()
}
