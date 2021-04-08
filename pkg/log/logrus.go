package log

import (
	"github.com/sirupsen/logrus"
)

func NewLogrusLogger(logger logrus.FieldLogger) *LogrusLogger {
	a := &AbstractLogger{}
	l := &LogrusLogger{
		a,
		logger.WithFields(logrus.Fields(Fields{})),
	}
	a.Logger = l
	return l
}

type LogrusLogger struct {
	*AbstractLogger
	entry *logrus.Entry
}

func (s *LogrusLogger) WithFields(fields Fields) Logger {
	return NewLogrusLogger(s.entry.WithFields(logrus.Fields(fields)))
}

func (s *LogrusLogger) Log(level Level, args ...interface{}) {
	logrusLevel, err := s.convertLevel(level)
	if err != nil {
		s.entry.Errorf("logrus logger: logging: converting level %s: %v", level, err)
		logrusLevel = logrus.ErrorLevel
	}
	s.entry.Log(logrusLevel, args...)
}

func (s *LogrusLogger) convertLevel(level Level) (logrus.Level, error) {
	return logrus.ParseLevel(level.String())
}
