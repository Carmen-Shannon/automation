package worker

import (
	"fmt"
	"sync"
)

func NewWorker(id int, taskChan chan Task, doneChan chan int, stopChan chan int, errChan chan error) Worker {
	return &worker{
		mu:       sync.Mutex{},
		id:       id,
		taskChan: taskChan,
		doneChan: doneChan,
		stopChan: stopChan,
		errChan:  errChan,
		active:   false,
	}
}

type worker struct {
	mu sync.Mutex

	id     int
	active bool

	taskChan chan Task
	doneChan chan int
	stopChan chan int
	errChan  chan error
}

type Worker interface {
	IsActive() bool
	Start()
	Stop()
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
			case t, ok := <-w.taskChan:
				if !ok {
					w.sendError(fmt.Errorf("Worker %d detected Task channel closing, stopping", w.id))
					return
				}

				_, err := t.Do()
				if err != nil {
					w.sendError(fmt.Errorf("Worker %d error: %w", w.id, err))
				}

				w.done()
			case i, ok := <-w.stopChan:
				if !ok || i == w.id {
					w.sendError(fmt.Errorf("Worker %d detected stop channel closing, stopping", w.id))
					return
				}
			}
		}
	}()
}

func (w *worker) Stop() {
	w.stopChan <- w.id
	w.done()
	w.mu.Lock()
	w.active = false
	w.mu.Unlock()
}

func (w *worker) done() {
	w.doneChan <- w.id
}

func (w *worker) sendError(err error) {
	w.errChan <- err
}
