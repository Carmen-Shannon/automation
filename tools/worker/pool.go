package worker

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

type dynamicWorkerPool struct {
	mu   sync.Mutex
	cond sync.Cond

	poolCtx    context.Context
	poolCancel context.CancelFunc

	workers []Worker

	taskQueue     chan Task
	stopChan      chan int
	errChan       chan error
	maxWorkers    int
	activeWorkers int
	stopped       bool
}

type DynamicWorkerPool interface {
	// ClearTaskQueue clears the task queue without stopping the workers.
	// This is useful for resetting the pool state without terminating the workers.
	//
	// It does not block the caller and returns immediately after clearing the queue.
	// Note: This method does not stop the workers, it only clears the task queue.
	// If you want to stop the workers, use the StopAll method instead.
	ClearTaskQueue()

	// DecreaseMaxWorkers decreases the maximum number of workers in the pool.
	// It does not block the caller and returns immediately after decreasing the number of workers.
	//
	// Note: This method will stop active workers if there are no inactive workers to remove.
	//
	// Parameters:
	//   - n: The number of workers to remove from the pool, must be less than or equal to the current number of workers.
	DecreaseMaxWorkers(n int)

	// GetMaxWorkers returns the maximum number of workers in the pool.
	//
	// Returns:
	//   - int: The maximum number of workers in the pool.
	GetMaxWorkers() int

	// IncreaseMaxWorkers increases the maximum number of workers in the pool.
	// It does not block the caller and returns immediately after increasing the number of workers.
	// The new workers will be initialized and added to the pool.
	//
	// Parameters:
	//   - n: The number of new workers to add to the pool.
	IncreaseMaxWorkers(n int)

	// IsWorking checks if the pool is currently processing tasks.
	// It returns true if there are tasks in the queue or if any workers are active.
	// This method is non-blocking and returns immediately.
	//
	// Note: this method should not be looped on, as it may cause a busy wait.
	// Instead, use the Wait method to block until all tasks are completed.
	IsWorking() bool

	// Stop stops all workers in the pool.
	// It does not clear the task queue, so any tasks that are currently in the queue will remain there and be picked up by the scheduler.
	Stop()

	// Start re-starts the task handler so workers can be assigned tasks.
	Start()

	// SubmitTask submits a task to the pool for processing.
	// It does not block the caller and returns immediately after submitting the task.
	// The task will be processed by one of the available workers in the pool.
	//
	// Parameters:
	//   - t: The task to be submitted.
	SubmitTask(t Task)

	// Wait blocks until all tasks in the queue are completed and all workers are idle.
	// It is a blocking call and will not return until all tasks are processed.
	// This method is useful for waiting for all tasks to complete before proceeding with the next steps in your program.
	Wait()
}

var _ DynamicWorkerPool = (*dynamicWorkerPool)(nil)

func NewDynamicWorkerPool(maxWorkers int, queueSize int) DynamicWorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	pool := &dynamicWorkerPool{
		mu:         sync.Mutex{},
		taskQueue:  make(chan Task, queueSize),
		stopChan:   make(chan int, maxWorkers),
		errChan:    make(chan error, maxWorkers),
		maxWorkers: maxWorkers,
	}
	pool.cond = sync.Cond{L: &pool.mu}
	pool.poolCtx, pool.poolCancel = context.WithCancel(context.Background())

	pool.initWorkers()

	go pool.errorHandler()
	go pool.taskHandler()

	return pool
}

func (p *dynamicWorkerPool) ClearTaskQueue() {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Drain the taskQueue channel
	for len(p.taskQueue) > 0 {
		<-p.taskQueue
	}
}

func (p *dynamicWorkerPool) Stop() {
	p.poolCancel()
	for _, worker := range p.workers {
		if worker.IsActive() {
			worker.Stop()
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.activeWorkers = 0
	p.stopped = true
}

func (p *dynamicWorkerPool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.activeWorkers == 0 {
		go p.taskHandler()
		p.stopped = false
	}
}

func (p *dynamicWorkerPool) SubmitTask(t Task) {
	p.taskQueue <- t
}

func (p *dynamicWorkerPool) DecreaseMaxWorkers(n int) {
	if n > p.maxWorkers {
		n = p.maxWorkers
	}

	removed := 0

	p.mu.Lock()
	defer p.mu.Unlock()
	for i, w := range p.workers {
		if !w.IsActive() {
			w.Stop()
			p.workers = slices.Delete(p.workers, i, i+1)
			removed++
			if removed >= n {
				return
			}
		}
	}

	// if we removed all inactive workers and still need to remove more, stop active workers
	for i, w := range p.workers {
		if w.IsActive() {
			w.Stop()
			p.workers = slices.Delete(p.workers, i, i+1)
			removed++
			if removed >= n {
				return
			}
		}
	}
}

func (p *dynamicWorkerPool) GetMaxWorkers() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.maxWorkers
}

func (p *dynamicWorkerPool) IncreaseMaxWorkers(n int) {
	if n <= 0 {
		return
	}
	p.maxWorkers += n
	for range n {
		p.addWorker()
	}
}

func (p *dynamicWorkerPool) IsWorking() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopped || len(p.taskQueue) > 0 || p.activeWorkers > 0
}

func (p *dynamicWorkerPool) Wait() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.taskQueue) > 0 || p.activeWorkers > 0 {
		p.cond.Wait()
	}
}

func (p *dynamicWorkerPool) addWorker() {
	if len(p.workers) < p.maxWorkers {
		worker := NewWorker(len(p.workers), p.taskQueue, p.stopChan, p.errChan)
		p.mu.Lock()
		p.workers = append(p.workers, worker)
		p.mu.Unlock()
	}
}

func (p *dynamicWorkerPool) initWorkers() {
	for i := range p.maxWorkers {
		worker := NewWorker(i, p.taskQueue, p.stopChan, p.errChan)
		p.mu.Lock()
		p.workers = append(p.workers, worker)
		p.mu.Unlock()
	}
}

func (p *dynamicWorkerPool) taskHandler() {
	for {
		select {
		case <-p.poolCtx.Done():
			return
		default:
			queueSize := len(p.taskQueue)
			if queueSize > 0 && p.activeWorkers < p.maxWorkers {
				for _, worker := range p.workers {
					if !worker.IsActive() {
						worker.Start()

						// Lock only when modifying shared state
						p.mu.Lock()
						p.activeWorkers++
						p.mu.Unlock()

						if p.activeWorkers >= p.maxWorkers {
							break
						}
					}
				}
			}

			if queueSize == 0 && p.activeWorkers > 0 {
				for _, worker := range p.workers {
					if worker.IsActive() {
						worker.Stop()

						// Lock only when modifying shared state
						p.mu.Lock()
						p.activeWorkers--
						// Signal the condition variable if the queue is empty and no workers are active
						if len(p.taskQueue) == 0 && p.activeWorkers == 0 {
							p.cond.Signal()
						}
						p.mu.Unlock()

						if p.activeWorkers == 0 {
							break
						}
					}
				}
			}
		}
	}
}

func (p *dynamicWorkerPool) errorHandler() {
	for {
		select {
		case e, ok := <-p.errChan:
			if !ok {
				return
			}
			fmt.Println("Error detected in worker pool: ", e)
		default:
			return
		}
	}
}
