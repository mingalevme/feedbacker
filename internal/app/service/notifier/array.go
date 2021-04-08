package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
)

type ArrayNotifier struct {
	storage []model.Feedback
	logger  log.Logger
}

// Sync
func (s *ArrayNotifier) Notify(f model.Feedback) error {
	s.logger.WithField("_notifier", fmt.Sprintf("%T", s)).Infof("Notifying:\n%s", feedbackToMessage(f, &indent))
	return nil
}

func NewArrayNotifier(logger log.Logger) *ArrayNotifier {
	if logger == nil {
		panic("logger is nil")
	}
	return &ArrayNotifier{
		storage: []model.Feedback{},
		logger: logger,
	}
}
