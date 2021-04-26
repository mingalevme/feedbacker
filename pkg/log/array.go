package log

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

var storages = map[int][]string{}
var arrLogMutex sync.Mutex

type ArrayLogger struct {
	*AbstractLogger
	level  Level
	id     int
	fields Fields
	err    error
	req    *http.Request
}

func NewArrayLogger(level Level) *ArrayLogger {
	arrLogMutex.Lock()
	var id int
	for {
		id = rand.Int()
		if _, ok := storages[id]; !ok {
			storages[id] = []string{}
			arrLogMutex.Unlock()
			break
		}
	}
	a := &AbstractLogger{}
	l := &ArrayLogger{
		AbstractLogger: a,
		level:          level,
		id:             id,
		fields:         Fields{},
		err:            nil,
		req:            nil,
	}
	a.Logger = l
	return l
}

func (s *ArrayLogger) Clone() *ArrayLogger {
	clone := NewArrayLogger(s.level)
	clone.id = s.id
	clone.fields = s.fields.Clone()
	clone.err = s.err
	clone.req = s.req
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
	if level.isLower(s.level) {
		return
	}
	args = append(args, map[string]interface{}{
		"level":   level.String(),
		"fields":  s.fields,
		"request": s.req,
		"error":   s.err,
	})
	message := fmt.Sprint(args...)
	arrLogMutex.Lock()
	defer arrLogMutex.Unlock()
	storages[s.id] = append(storages[s.id], message)
}

func (s *ArrayLogger) Storage() []string {
	arrLogMutex.Lock()
	defer arrLogMutex.Unlock()
	return storages[s.id]
}
