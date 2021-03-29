package di

import (
	"database/sql"
	"github.com/mingalevme/feedbacker/app/service"
	"github.com/mingalevme/feedbacker/domain/feedback"
	"github.com/mingalevme/feedbacker/infrastructure/cfg"
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/mingalevme/feedbacker/infrastructure/log"
	"github.com/mingalevme/feedbacker/infrastructure/persistence"
	"github.com/mingalevme/feedbacker/infrastructure/persistence/db"
	"github.com/mingalevme/feedbacker/infrastructure/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Container interface {
	GetFeedbackRepository() feedback.Repository
	GetLogger() log.Logger
	GetLeaveFeedbackService() service.LeaveFeedbackService
}

type container struct {
	config               cfg.Config
	logger               log.Logger
	feedbackRepository   feedback.Repository
	db                   *sql.DB
	leaveFeedbackService service.LeaveFeedbackService
}

func New(config cfg.Config) Container {
	return &container{
		config: config,
	}
}

func (s *container) GetLogger() log.Logger {
	if s.logger == nil {
		logger := logrus.New()
		// @TODO: parse environment
		logger.SetOutput(os.Stdout)
		if level, err := logrus.ParseLevel(s.config.GetEnvVar("LOG_LEVEL", "debug")); err != nil {
			panic(errors.Wrap(err, "Error while parsing log level"))
		} else {
			logger.SetLevel(level)
		}
		s.logger = logger
	}
	return s.logger
}

func (s *container) GetFeedbackRepository() feedback.Repository {
	driver := s.config.GetEnvVar("PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		connection, err := s.getDatabaseConnection()
		if err != nil {
			panic(errors.Wrap(err, "Error while initializing connection to database"))
		}
		return db.NewFeedbackRepository(connection)
	} else if driver == "logger" {
		return persistence.NewLoggerFeedbackRepository(s.GetLogger())
	}
	panic(errors.Errorf("Unsupported persistence driver: %s", driver))
}

func (s *container) GetLeaveFeedbackService() service.LeaveFeedbackService {
	if s.leaveFeedbackService == nil {
		s.leaveFeedbackService = service.NewLeaveFeedbackService(s.GetFeedbackRepository(), s.GetLogger())
	}
	return s.leaveFeedbackService
}

func (s *container) getDatabaseConnection() (*sql.DB, error) {
	if s.db != nil {
		return s.db, nil
	}
	params := map[string]interface{}{
		"Host":     s.config.GetEnvVar("DB_HOST", "127.0.0.1"),
		"Port":     s.config.GetEnvVar("DB_PORT", "5432"),
		"User":     s.config.GetEnvVar("DB_USER", "postgres"),
		"Pass":     s.config.GetEnvVar("DB_PASSWORD", "postgres"),
		"Database": s.config.GetEnvVar("DB_NAME", "postgres"),
	}
	dataSourceName := util.Sprintf("postgres://%{User}s:%{Pass}s@%{Host}s:%{Port}s/%{Database}s?sslmode=disable", params)
	connection, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "Error while initializing connection to database")
	}
	if err = connection.Ping(); err != nil {
		return nil, errors.Wrap(err, "Error while pinging connection to database")
	}
	s.db = connection
	return s.db, nil
}

var instance Container

func GetDefault() Container {
	if instance == nil {
		instance = New(cfg.New(env.New()))
	}
	return instance
}
