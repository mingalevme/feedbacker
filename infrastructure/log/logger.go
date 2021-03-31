package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Fields logrus.Fields

type Logger interface {
	WithField(key string, value interface{}) Entry
	WithFields(fields Fields) Entry
	WithError(err error) Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	// IsDebugEnabled() bool
	// IsInfoEnabled() bool
	// IsWarnEnabled() bool
	// IsErrorEnabled() bool
	// IsFatalEnabled() bool
	// IsPanicEnabled() bool
}

type Entry interface {
	Logger
	WithContext(ctx context.Context) Entry
	WithTime(t time.Time) Entry
	String() (string, error)
}
