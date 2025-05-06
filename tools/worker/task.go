package worker

type Task struct {
	ID      int
	Payload any
	Do      func() (any, error)
}
