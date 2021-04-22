package log

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

// ZapLogger is just an example of wrapping the zap-logger in the app Logger,
// but it breaks all it's (zap) benefits - zero allocations.
// So use the wrapper for educational purposes only.
type ZapLogger struct {
	*AbstractLogger
	zap *zap.Logger
	//level   Level
	fields Fields
	req    *http.Request
	err    error
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	a := &AbstractLogger{}
	l := &ZapLogger{
		AbstractLogger: a,
		zap:            logger,
		fields:         Fields{},
		req:            nil,
		err:            nil,
	}
	a.Logger = l
	return l
}

func (s *ZapLogger) Clone() *ZapLogger {
	clone := NewZapLogger(s.zap)
	clone.fields = s.fields.Clone()
	clone.req = s.req
	clone.err = s.err
	return clone
}

func (s *ZapLogger) WithFields(fields Fields) Logger {
	clone := s.Clone()
	for k, v := range fields {
		clone.fields[k] = v
	}
	return clone
}

func (s *ZapLogger) WithError(err error) Logger {
	clone := s.Clone()
	clone.err = err
	return clone
}

func (s *ZapLogger) WithRequest(req *http.Request) Logger {
	clone := s.Clone()
	clone.req = req
	return clone
}

func (s *ZapLogger) Log(level Level, args ...interface{}) {
	var size = len(s.fields)
	if s.err != nil {
		size++
	}
	if s.req != nil {
		size++
	}
	fields := make([]zap.Field, size)
	i := -1
	for k, v := range s.fields {
		i++
		fields[i] = zap.Any(k, v)
	}
	if s.err != nil {
		i++
		fields[i] = zap.Error(s.err)
	}
	if s.req != nil {
		i++
		fields[i] = zap.Reflect("request", RequestToMapTransformer(s.req))
	}
	message := fmt.Sprint(args...)
	switch level {
	case LevelDebug:
		s.zap.Debug(message, fields...)
	case LevelInfo:
		s.zap.Info(message, fields...)
	case LevelWarning:
		s.zap.Warn(message, fields...)
	case LevelError:
		s.zap.Error(message, fields...)
	case LevelFatal:
		s.zap.Fatal(message, fields...)
	default:
		s.zap.Fatal(fmt.Sprintf("%s: %s", "Invalid log level while logging to Zap-logger", message), fields...)
	}
}

func (s *ZapLogger) Close() {
	_ = s.zap.Sync()
}

func ParseZapLevel(lvl string) (zapcore.Level, error) {
	switch lvl {
	case LevelDebug.String():
		return zap.DebugLevel, nil
	case LevelInfo.String():
		return zap.InfoLevel, nil
	case LevelWarning.String():
		return zap.WarnLevel, nil
	case LevelError.String():
		return zap.ErrorLevel, nil
	case LevelFatal.String():
	case "panic":
	case "dpanic":
		return zap.FatalLevel, nil
	}
	return zap.DebugLevel, errors.Errorf("parsing zap lo level: invalid level: %s", lvl)
}
