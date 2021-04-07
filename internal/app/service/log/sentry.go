package log

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
	"reflect"
)

const sentryMaxErrorDepth = 10

type SentryLogger struct {
	*AbstractLogger
	hub *sentry.Hub
	level Level
	err error
}

func NewSentryLogger(hub *sentry.Hub, level Level) *SentryLogger {
	a := &AbstractLogger{}
	s := &SentryLogger{
		a,
		hub,
		level,
		nil,
	}
	a.Logger = s
	return s
}

func (s *SentryLogger) Clone() *SentryLogger {
	h := s.hub.Clone()
	clone := NewSentryLogger(h, s.level)
	clone.err = s.err
	return clone
}

//func (s *SentryLogger) WithField(key string, value interface{}) Logger {
//	return s.WithFields(Fields{key: value})
//}

func (s *SentryLogger) WithFields(fields Fields) Logger {
	clone := s.Clone()
	clone.hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtras(fields)
	})
	return clone
}

func (s *SentryLogger) WithError(err error) Logger {
	clone := s.Clone()
	clone.err = err
	return clone
}

func (s *SentryLogger) WithRequest(req *http.Request) Logger {
	clone := s.Clone()
	clone.hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetRequest(req)
	})
	return clone
}

func (s *SentryLogger) Log(level Level, args ...interface{}) {
	if !level.isGTE(s.level) {
		return
	}
	clone := s.Clone()
	clone.hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.Level(level))
		if clone.err != nil {
			scope.SetExtra("error", clone.err.Error())
			//scope.AddEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			//	event.Exception = s.errorToExceptions(clone.err)
			//	return event
			//})
		}
	})
	clone.hub.CaptureMessage(fmt.Sprint(args...))
}

func (s *SentryLogger) errorToExceptions(err error) []sentry.Exception {
	exceptions := []sentry.Exception{}

	if err == nil {
		return exceptions
	}

	for i := 0; i < sentryMaxErrorDepth && err != nil; i++ {
		exceptions = append(exceptions, sentry.Exception{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: sentry.ExtractStacktrace(err),
		})
		switch previous := err.(type) {
		case interface{ Unwrap() error }:
			err = previous.Unwrap()
		case interface{ Cause() error }:
			err = previous.Cause()
		default:
			err = nil
		}
	}

	// Add a trace of the current stack to the most recent error in a chain if
	// it doesn't have a stack trace yet.
	// We only add to the most recent error to avoid duplication and because the
	// current stack is most likely unrelated to errors deeper in the chain.
	if exceptions[0].Stacktrace == nil {
		exceptions[0].Stacktrace = sentry.NewStacktrace()
	}

	// event.Exception should be sorted such that the most recent error is last.
	for i := len(exceptions)/2 - 1; i >= 0; i-- {
		opp := len(exceptions) - 1 - i
		exceptions[i], exceptions[opp] = exceptions[opp], exceptions[i]
	}

	return exceptions
}

