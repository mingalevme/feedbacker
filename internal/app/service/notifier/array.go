package notifier

import (
	"fmt"
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
)

type ArrayFeedbackLeftNotifier struct {
	storage []model.Feedback
	logger  log.Logger
}

// Sync
func (s *ArrayFeedbackLeftNotifier) Notify(f model.Feedback) error {
	s.logger.WithField("_notifier", fmt.Sprintf("%T", s)).Infof("Notifying:\n%s", feedbackToMessage(f, &indent))
	return nil
}

func NewArrayFeedbackLeftNotifier(logger log.Logger) *ArrayFeedbackLeftNotifier {
	if logger == nil {
		panic("logger is nil")
	}
	return &ArrayFeedbackLeftNotifier{
		storage: []model.Feedback{},
		logger: logger,
	}
}
