package notifier

import (
	"github.com/mingalevme/feedbacker/internal/app/model"
	"github.com/mingalevme/feedbacker/pkg/log"
)

type LogNotifier struct {
	logger log.Logger
	level log.Level
}

func NewLogNotifier(logger log.Logger, level log.Level) *LogNotifier {
	return &LogNotifier{
		logger: logger,
		level: level,
	}
}

func (s *LogNotifier) Name() string {
	return "log"
}

func (s *LogNotifier) Health() error {
	return nil
}

func (s *LogNotifier) Notify(f model.Feedback) error {
	s.logger.Infof("Notifying:\n%s", feedbackToMessage(f, &indent))
	return nil
}
