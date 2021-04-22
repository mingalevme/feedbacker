package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
)

// ZerologLogger is just an example of wrapping the zerolog in the app Logger,
// but it breaks all it's (zerolog) benefits - zero allocations.
// So use the wrapper for educational purposes only.
type ZerologLogger struct {
	*AbstractLogger
	zlogger zerolog.Logger
	level   Level
	fields  Fields
	req     *http.Request
	err     error
}

func NewZerologLogger(logger zerolog.Logger) *ZerologLogger {
	a := &AbstractLogger{}
	l := &ZerologLogger{
		AbstractLogger: a,
		zlogger:        logger,
		fields:         Fields{},
		req:            nil,
		err:            nil,
	}
	a.Logger = l
	return l
}

func (s *ZerologLogger) Clone() *ZerologLogger {
	clone := NewZerologLogger(s.zlogger)
	clone.fields = s.fields.Clone()
	clone.req = s.req
	clone.err = s.err
	return clone
}

func (s *ZerologLogger) WithFields(fields Fields) Logger {
	clone := s.Clone()
	for k, v := range fields {
		clone.fields[k] = v
	}
	return clone
}

func (s *ZerologLogger) WithError(err error) Logger {
	clone := s.Clone()
	clone.err = err
	return clone
}

func (s *ZerologLogger) WithRequest(req *http.Request) Logger {
	clone := s.Clone()
	clone.req = req
	return clone
}

func (s *ZerologLogger) Log(level Level, args ...interface{}) {
	logger := s.zlogger
	if s.req != nil {
		logger = logger.With().Interface("request", RequestToMapTransformer(s.req)).Logger()
	}
	zerologLevel, err := zerolog.ParseLevel(level.String())
	if err != nil {
		s.zlogger.Error().Err(err).Msgf("Error while converting app log level (%s) to zerolog level", level)
		zerologLevel = zerolog.ErrorLevel
	}
	e := logger.WithLevel(zerologLevel)
	if s.fields != nil && len(s.fields) > 0 {
		e = e.Fields(s.fields)
	}
	if s.err != nil {
		e = e.Err(s.err)
	}
	e.Msg(fmt.Sprint(args...))
}
