package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"sync"
)

type ArrayNotifier struct {
	Storage []model.Feedback
	mu      sync.Mutex
}

func (s *ArrayNotifier) Name() string {
	return "array"
}

func (s *ArrayNotifier) Health() error {
	return nil
}

func (s *ArrayNotifier) Notify(f model.Feedback) error {
	s.mu.Lock()
	s.Storage = append(s.Storage, f)
	s.mu.Unlock()
	return nil
}

func NewArrayNotifier() *ArrayNotifier {
	n := &ArrayNotifier{
		Storage: []model.Feedback{},
	}
	return n
}
