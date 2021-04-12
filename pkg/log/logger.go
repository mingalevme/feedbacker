package log

import (
	"net/http"
)

type Fields map[string]interface{}

func (f Fields) Clone() Fields {
	if f == nil {
		return nil
	}
	clone := Fields{}
	for key, value := range f {
		clone[key] = value
	}
	return clone
}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	WithRequest(req *http.Request) Logger

	Log(level Level, args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	//Printf(format string, args ...interface{})
	//Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	//Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	//Print(args ...interface{})
	//Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	//Panic(args ...interface{})

	//Debugln(args ...interface{})
	//Infoln(args ...interface{})
	//Println(args ...interface{})
	//Warnln(args ...interface{})
	//Warningln(args ...interface{})
	//Errorln(args ...interface{})
	//Fatalln(args ...interface{})
	//Panicln(args ...interface{})
}

//type Entry interface {
//	Logger
//	WithContext(ctx context.Context) Entry
//	WithTime(t time.Time) Entry
//	String() (string, error)
//}

var RequestTransformer = func(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"url":     r.URL.String(),
		"method":  r.Method,
		"headers": r.Header,
	}
}
