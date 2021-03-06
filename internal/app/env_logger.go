package app

import (
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"os"
	"strings"
)

func (s *Container) Logger() log.Logger {
	if s.logger != nil {
		return s.logger
	}
	channel := s.EnvVarBag.Get("LOG_CHANNEL", "stdout")
	s.logger = s.newLogChannel(channel)
	return s.logger
}

// todo: go-kit, apex/log
func (s *Container) newLogChannel(channel string) log.Logger {
	switch channel {
	case "stdout":
		return s.newStdoutLogger()
	case "stderr":
		return s.newStderrLogger()
	case "zerolog":
		return s.newZerologLogger()
	case "zap":
		return s.newZapLogger()
	case "sentry":
		return s.newSentryLogger()
	case "rollbar":
		return s.newRollbarLogger()
	case "stack":
		return s.newStackLogger()
	case "array":
		if lvl, err := log.ParseLevel(s.EnvVarBag.Get("LOG_ARRAY_LEVEL", "debug")); err != nil {
			panic(errors.Wrap(err, "parsing LOG_ARRAY_LEVEL-envvar"))
		} else {
			return log.NewArrayLogger(lvl)
		}
	case "null":
		return s.newNullLogger()
	default:
		panic(errors.Errorf("unsupported log channel: %s", channel))
	}
}

func (s *Container) newStackLogger() log.Logger {
	logger := log.NewStackLogger()
	channels := strings.Split(s.EnvVarBag.Get("LOG_STACK_CHANNELS", "stdout"), ",")
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

func (s *Container) newNullLogger() log.Logger {
	return log.NewNullLogger()
}

func (s *Container) newStdoutLogger() log.Logger {
	if level, err := log.ParseLevel(s.EnvVarBag.Get("LOG_STDOUT_LEVEL", "debug")); err != nil {
		panic(errors.Wrap(err, "parsing stdout logging level"))
	} else {
		return log.NewStdoutLogger(level)
	}
}

func (s *Container) newStderrLogger() log.Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stderr)
	if level, err := logrus.ParseLevel(s.EnvVarBag.Get("LOG_STDERR_LEVEL", "error")); err != nil {
		panic(errors.Wrap(err, "parsing stderr logging level"))
	} else {
		logrusLogger.SetLevel(level)
	}
	return log.NewLogrusLogger(logrusLogger)
}

func (s *Container) newZerologLogger() log.Logger {
	z := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if level, err := zerolog.ParseLevel(s.EnvVarBag.Get("LOG_ZEROLOG_LEVEL", "debug")); err != nil {
		panic(errors.Wrap(err, "parsing zerolog logging level"))
	} else {
		z = z.Level(level)
	}
	return log.NewZerologLogger(z)
}

func (s *Container) newSentryLogger() log.Logger {
	level, err := log.ParseLevel(s.EnvVarBag.Get("LOG_SENTRY_LEVEL", log.LevelWarning.String()))
	if err != nil {
		panic(errors.Wrap(err, "parsing sentry log level (LOG_SENTRY_LEVEL)"))
	}
	return log.NewSentryLogger(s.Sentry(), level)
}

func (s *Container) newRollbarLogger() log.Logger {
	level, err := log.ParseLevel(s.EnvVarBag.Get("LOG_ROLLBAR_LEVEL", log.LevelWarning.String()))
	if err != nil {
		panic(errors.Wrap(err, "parsing rollbar log level (LOG_ROLLBAR_LEVEL)"))
	}
	return log.NewRollbarLogger(s.Rollbar(), level)
}

func (s *Container) newZapLogger() log.Logger {
	cfg := zap.NewProductionConfig()
	if zapLevel, err := log.ParseZapLevel(s.EnvVarBag.Get("LOG_ZAP_LEVEL", "debug")); err != nil {
		panic(errors.Wrap(err, "parsing zap logging level"))
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	}
	if z, err := cfg.Build(); err != nil {
		panic(errors.Wrap(err, "building zap logger"))
	} else {
		return log.NewZapLogger(z)
	}
}
