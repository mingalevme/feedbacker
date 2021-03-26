package infrastructure

import (
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New()
		logger.SetOutput(os.Stdout)
		if level, err := log.ParseLevel(env.GetEnvValue("FEEDBACKER_LOG_LEVEL", "debug")); err != nil {
			panic(errors.Wrap(err, "Error while parsing log level"))
		} else {
			logger.SetLevel(level)
		}
	}
	return logger
}
