package di

import (
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (s *container) GetLogger() log.Logger {
	if s.logger != nil {
		return s.logger
	}
	channel := s.config.GetEnvVar("LOG_CHANNEL", "stdout")
	s.logger = s.newLogChannel(channel)
	return s.logger
}

func (s *container) newLogChannel(channel string) log.Logger {
	if channel == "null" {
		return s.newNullLogger()
	}
	if channel == "stdout" {
		return s.newStdoutLogger()
	}
	if channel == "sentry" {
		return s.newSentryLogger()
	}
	if channel == "rollbar" {
		return s.newRollbarLogger()
	}
	if channel == "stack" {
		return s.newStackLogger()
	}
	panic(errors.Errorf("unsupported log channel: %s", channel))
}

func (s *container) newStackLogger() log.Logger {
	logger := log.NewStackLogger()
	channels := strings.Split(s.config.GetEnvVar("LOG_STACK_CHANNELS", "stdout"), ",")
	for _, channel := range channels {
		channel = strings.TrimSpace(channel)
		if channel == "" {
			continue
		}
		if channel == "stack" {
			panic(errors.Errorf("stack channel recursion"))
		}
		logger.Add(s.newLogChannel(channel))
	}
	return logger
}

func (s *container) newNullLogger() log.Logger {
	return log.NewNullLogger()
}

func (s *container) newStdoutLogger() log.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stdout)
	if level, err := logrus.ParseLevel(s.config.GetEnvVar("LOG_STDOUT_LEVEL", "debug")); err != nil {
		panic(errors.Wrap(err, "parsing stdout logging level"))
	} else {
		logrusLogger.SetLevel(level)
	}
	return log.NewLogrusLogger(logrusLogger)
}

func (s *container) newSentryLogger() log.Logger {
	level, err := log.ParseLevel(s.config.GetEnvVar("LOG_SENTRY_LEVEL", log.LevelWarning.String()))
	if err != nil {
		panic(errors.Wrap(err, "parsing sentry log level (LOG_SENTRY_LEVEL)"))
	}
	return log.NewSentryLogger(s.GetSentry(), level)
}

func (s *container) newRollbarLogger() log.Logger {
	level, err := log.ParseLevel(s.config.GetEnvVar("LOG_ROLLBAR_LEVEL", log.LevelWarning.String()))
	if err != nil {
		panic(errors.Wrap(err, "parsing rollbar log level (LOG_ROLLBAR_LEVEL)"))
	}
	return log.NewRollbarLogger(s.GetRollbar(), level)
}
