package log

import "net/http"

type StackLogger struct {
	*AbstractLogger
	loggers []Logger
}

func NewStackLogger() *StackLogger {
	a := &AbstractLogger{}
	s := &StackLogger{
		a,
		[]Logger{},
	}
	a.Logger = s
	return s
}

func (s *StackLogger) WithField(key string, value interface{}) Logger {
	stack := NewStackLogger()
	for _, logger := range s.loggers {
		stack.loggers = append(stack.loggers, logger.WithField(key, value))
	}
	return stack
}

func (s *StackLogger) WithFields(fields Fields) Logger {
	stack := NewStackLogger()
	for _, logger := range s.loggers {
		stack.loggers = append(stack.loggers, logger.WithFields(fields))
	}
	return stack
}

func (s *StackLogger) WithError(err error) Logger {
	stack := NewStackLogger()
	for _, logger := range s.loggers {
		stack.loggers = append(stack.loggers, logger.WithError(err))
	}
	return stack
}

func (s *StackLogger) WithRequest(req *http.Request) Logger {
	stack := NewStackLogger()
	for _, logger := range s.loggers {
		stack.loggers = append(stack.loggers, logger.WithRequest(req))
	}
	return stack
}

func (s *StackLogger) Log(level Level, args ...interface{}) {
	for _, logger := range s.loggers {
		logger.Log(level, args...)
	}
}

func (s *StackLogger) Add(logger Logger) {
	s.loggers = append(s.loggers, logger)
}

func (s *StackLogger) Close() {
	for _, logger := range s.loggers {
		logger.Close()
	}
}
