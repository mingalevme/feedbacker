package log

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Logger interface {
	logrus.FieldLogger
}

func NewNullLogger() Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}
