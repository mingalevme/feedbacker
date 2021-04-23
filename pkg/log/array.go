package log

import (
	"fmt"
	"net/http"
)

type ArrayLogger struct {
	*AbstractLogger
	fields  Fields
	err     error
	req     *http.Request
	Storage []string
}

func NewArrayLogger() *ArrayLogger {
	a := &AbstractLogger{}
	l := &ArrayLogger{
		AbstractLogger: a,
		fields:         Fields{},
		err:            nil,
		req:            nil,
		Storage:        []string{},
	}
	a.Logger = l
	return l
}

func (s *ArrayLogger) Clone() *ArrayLogger {
	clone := NewArrayLogger()
	clone.fields = s.fields.Clone()
	clone.err = s.err
	clone.req = s.req
	clone.Storage = s.Storage
	return clone
}

func (s *ArrayLogger) WithFields(fields Fields) Logger {
	clone := s.Clone()
	clone.fields = fields
	return clone
}

func (s *ArrayLogger) WithError(err error) Logger {
	clone := s.Clone()
	clone.err = err
	return clone
}

func (s *ArrayLogger) WithRequest(req *http.Request) Logger {
	clone := s.Clone()
	clone.req = req
	return clone
}

func (s *ArrayLogger) Log(level Level, args ...interface{}) {
	args = append(args, map[string]interface{}{
		"level": level.String(),
		"fields": s.fields,
		"request": s.req,
		"error": s.err,
	})
	message := fmt.Sprint(args...)
	s.Storage = append(s.Storage, message)
}
