package log

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Level string

const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelFatal   Level = "fatal"
)

var levelMap = map[Level]int{
	LevelDebug:   100,
	LevelInfo:    200,
	LevelWarning: 300,
	LevelError:   400,
	LevelFatal:   500,
}

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarning, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	}
	return Level(""), fmt.Errorf("not a valid log level: %q", lvl)
}

func (l Level) String() string {
	return string(l)
}

func (l Level) isGTE(t Level) bool {
	v1, ok := levelMap[l]
	if ok == false {
		panic(errors.Errorf("invalid source level: %s", l))
	}
	v2, ok := levelMap[t]
	if ok == false {
		panic(errors.Errorf("invalid targer level: %s", t))
	}
	return v1 >= v2
}

func (l Level) isLower(t Level) bool {
	return !l.isGTE(t)
}
