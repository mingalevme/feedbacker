package log

import (
	"fmt"
	"net/http"
)

type AbstractLogger struct {
	Logger
}

func (s *AbstractLogger) WithField(key string, value interface{}) Logger {
	return s.WithFields(Fields{key: value})
}

//func (s *AbstractLogger) WithFields(fields Fields) Logger {
//	panic("not implemented")
//}

func (s *AbstractLogger) WithError(err error) Logger {
	return s.WithField("error", err)
	//panic("not implemented")
}

func (s *AbstractLogger) WithRequest(req *http.Request) Logger {
	return s.WithField("request", req)
}

//func (s *AbstractLogger) Log(level Level, args ...interface{}) {
//	panic("not implemented")
//}

func (s *AbstractLogger) Logf(level Level, format string, args ...interface{}) {
	s.Log(level, fmt.Sprintf(format, args...))
}

func (s *AbstractLogger) Debugf(format string, args ...interface{}) {
	s.Logf(LevelDebug, format, args...)
}

func (s *AbstractLogger) Infof(format string, args ...interface{}) {
	s.Logf(LevelInfo, format, args...)
}

func (s *AbstractLogger) Warningf(format string, args ...interface{}) {
	s.Logf(LevelWarning, format, args...)
}

func (s *AbstractLogger) Errorf(format string, args ...interface{}) {
	s.Logf(LevelError, format, args...)
}

func (s *AbstractLogger) Fatalf(format string, args ...interface{}) {
	s.Logf(LevelFatal, format, args...)
}

func (s *AbstractLogger) Debug(args ...interface{}) {
	s.Log(LevelDebug, args...)
}

func (s *AbstractLogger) Info(args ...interface{}) {
	s.Log(LevelInfo, args...)
}

func (s *AbstractLogger) Warning(args ...interface{}) {
	s.Log(LevelWarning, args...)
}

func (s *AbstractLogger) Error(args ...interface{}) {
	s.Log(LevelError, args...)
}

func (s *AbstractLogger) Fatal(args ...interface{}) {
	s.Log(LevelFatal, args...)
}

func Close() {}
