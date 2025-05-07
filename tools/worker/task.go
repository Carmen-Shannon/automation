package worker

// Task represents a task to be processed by a worker.
// TODO: turn this into an interface, set up an easier-to-use option builder pattern
type Task struct {
	ID      int
	Payload any
	Do      func() (any, error)
}
