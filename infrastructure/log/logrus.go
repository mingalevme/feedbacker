package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

func NewLogrusLogger(logrus logrus.FieldLogger) Logger {
	return &logrusLogger{
		FieldLogger: logrus,
	}
}

type logrusLogger struct {
	logrus.FieldLogger
}

func (s * logrusLogger) WithField(key string, value interface{}) Entry {
	return &entry{
		s.FieldLogger.WithField(key, value),
	}
}

func (s * logrusLogger) WithFields(fields Fields) Entry {
	return &entry{
		s.FieldLogger.WithFields(logrus.Fields(fields)),
	}
}

func (s * logrusLogger) WithError(err error) Entry {
	return &entry{
		s.FieldLogger.WithError(err),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type entry struct {
	*logrus.Entry
}

func (s *entry) WithField(key string, value interface{}) Entry {
	return &entry{
		s.Entry.WithField(key, value),
	}
}

func (s *entry) WithFields(fields Fields) Entry {
	return &entry{
		s.Entry.WithFields(logrus.Fields(fields)),
	}
}

func (s *entry) WithError(err error) Entry {
	return &entry{
		s.Entry.WithError(err),
	}
}

func (s *entry) WithContext(ctx context.Context) Entry {
	return &entry{
		s.Entry.WithContext(ctx),
	}
}

func (s *entry) WithTime(t time.Time) Entry {
	return &entry{
		s.Entry.WithTime(t),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func NewNullLogger() Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return NewLogrusLogger(logger)
}
