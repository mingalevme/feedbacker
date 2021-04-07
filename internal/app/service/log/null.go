package log

type NullLogger struct {
	*AbstractLogger
}

func NewNullLogger() *NullLogger {
	a := &AbstractLogger{}
	l := &NullLogger{
		a,
	}
	a.Logger = l
	return l
}

func (s *NullLogger) WithFields(fields Fields) Logger {
	return s
}

func (s *NullLogger) Log(level Level, args ...interface{}) {

}
