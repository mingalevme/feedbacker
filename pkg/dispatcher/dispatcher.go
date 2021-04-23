package dispatcher

import "github.com/pkg/errors"

var MaxQueueSizeReached = errors.New("maximum queue size has been reached")
var TaskQueueIsExiting = errors.New("task queue is exiting")

type Dispatcher interface {
	Name() string
	Enqueue(t Task) error
	Run() error
	Stop() error
	Health() error
}

type Task func() error

type TaskQueue interface {
	Enqueue(t Task) error
}

type taskQueue struct {
	dispatcher Dispatcher
}

func NewTaskQueue(d Dispatcher) TaskQueue {
	return &taskQueue{
		dispatcher: d,
	}
}

func (s *taskQueue) Enqueue(t Task) error {
	return s.dispatcher.Enqueue(t)
}
