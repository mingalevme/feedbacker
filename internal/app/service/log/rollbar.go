package log

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
	"net/http"
)

type RollbarLogger struct {
	*AbstractLogger
	rollbar *rollbar.Client
	level   Level
	fields  Fields
	err     error
	req     *http.Request
}

func NewRollbarLogger(client *rollbar.Client, level Level) *RollbarLogger {
	a := &AbstractLogger{}
	l := &RollbarLogger{
		AbstractLogger: a,
		rollbar:        client,
		level:          level,
		fields:         Fields{},
		err:            nil,
		req:            nil,
	}
	a.Logger = l
	return l
}

func (s *RollbarLogger) Clone() *RollbarLogger {
	clone := NewRollbarLogger(s.rollbar, s.level)
	clone.fields = s.fields.Clone()
	clone.err = s.err
	clone.req = s.req
	return clone
}

func (s *RollbarLogger) WithFields(fields Fields) Logger {
	clone := s.Clone()
	clone.fields = fields
	return clone
}

func (s *RollbarLogger) WithError(err error) Logger {
	clone := s.Clone()
	clone.err = err
	return clone
}

func (s *RollbarLogger) WithRequest(req *http.Request) Logger {
	clone := s.Clone()
	clone.req = req
	return clone
}

func (s *RollbarLogger) Log(level Level, args ...interface{}) {
	if level.isLower(s.level) {
		return
	}
	message := fmt.Sprint(args...)
	if s.err != nil && s.req != nil {
		s.rollbar.RequestErrorWithExtras(s.convertLevel(level), s.req, errors.Wrap(s.err, message), s.fields)
	} else if s.err != nil {
		s.rollbar.ErrorWithExtras(s.convertLevel(level), errors.Wrap(s.err, message), s.fields)
	} else if s.req != nil {
		s.rollbar.RequestMessageWithExtras(s.convertLevel(level), s.req, message, s.fields)
	} else {
		s.rollbar.MessageWithExtras(s.convertLevel(level), message, s.fields)
	}
}

func (s *RollbarLogger) convertLevel(level Level) string {
	l := level.String()
	if l == LevelFatal.String() {
		l = "critical"
	}
	return l
}
