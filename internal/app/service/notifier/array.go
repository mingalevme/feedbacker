package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
	"sync"
)

type ArrayNotifier struct {
	Storage []model.Feedback
	Logger  log.Logger
	mu      sync.Mutex
}

func (s *ArrayNotifier) Name() string {
	return "array"
}

func (s *ArrayNotifier) Health() error {
	return nil
}

// Sync
func (s *ArrayNotifier) Notify(f model.Feedback) error {
	s.mu.Lock()
	s.Storage = append(s.Storage, f)
	s.mu.Unlock()
	s.Logger.WithField("_notifier", fmt.Sprintf("%T", s)).Infof("Notifying:\n%s", feedbackToMessage(f, &indent))
	return nil
}

func NewArrayNotifier(logger log.Logger) *ArrayNotifier {
	n := &ArrayNotifier{
		Storage: []model.Feedback{},
		Logger:  log.NewNullLogger(),
	}
	if logger != nil {
		n.Logger = logger
	}
	return n
}
