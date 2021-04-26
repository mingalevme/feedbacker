package log

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func NewStdoutLogger(lvl Level) *LogrusLogger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	if level, err := logrus.ParseLevel(lvl.String()); err != nil {
		panic(errors.Wrap(err, "parsing app log level to logrus level"))
	} else {
		logger.SetLevel(level)
	}
	return NewLogrusLogger(logger)
}


