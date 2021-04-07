package di

import (
	"database/sql"
	"github.com/mingalevme/feedbacker/_ddd/app/service"
	"github.com/mingalevme/feedbacker/_ddd/domain/feedback"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/config"
	feedbackRepository "github.com/mingalevme/feedbacker/_ddd/infrastructure/persistence/repository/feedback"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/util"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Container interface {
	GetFeedbackRepository() feedback.Repository
	GetLogger() log.Logger
	GetLeaveFeedbackService() service.LeaveFeedbackService
	GetViewFeedbackService() service.ViewFeedbackService
}

type container struct {
	config               config.Config
	logger               log.Logger
	feedbackRepository   feedback.Repository
	db                   *sql.DB
	leaveFeedbackService service.LeaveFeedbackService
	viewFeedbackService  service.ViewFeedbackService
}

func New(config config.Config) Container {
	instance = &container{
		config: config,
	}
	return instance
}

func (s *container) GetLogger() log.Logger {
	if s.logger != nil {
		return s.logger
	}
	logrusLogger := logrus.New()
	// @TODO: parse environment
	logrusLogger.SetOutput(os.Stdout)
	if level, err := logrus.ParseLevel(s.config.GetEnvVar("LOG_LEVEL", "info")); err != nil {
		panic(errors.Wrap(err, "Error while parsing loggin level"))
	} else {
		logrusLogger.SetLevel(level)
	}
	s.logger = log.NewLogrusLogger(logrusLogger)
	return s.logger
}

func (s *container) GetFeedbackRepository() feedback.Repository {
	if s.feedbackRepository != nil {
		return s.feedbackRepository
	}
	driver := s.config.GetEnvVar("PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		s.feedbackRepository = feedbackRepository.NewFeedbackRepository(s.getDatabaseConnection(), s.GetLogger())
	} else if driver == "array" {
		s.feedbackRepository = feedbackRepository.NewArrayFeedbackRepository(s.GetLogger())
	} else {
		panic(errors.Errorf("Unsupported persistence driver: %s", driver))
	}
	return s.feedbackRepository
}

func (s *container) GetLeaveFeedbackService() service.LeaveFeedbackService {
	if s.leaveFeedbackService == nil {
		s.leaveFeedbackService = service.NewLeaveFeedbackService(s.GetFeedbackRepository(), s.GetLogger())
	}
	return s.leaveFeedbackService
}

func (s *container) GetViewFeedbackService() service.ViewFeedbackService {
	if s.viewFeedbackService == nil {
		s.viewFeedbackService = service.NewViewFeedbackService(s.GetFeedbackRepository(), s.GetLogger())
	}
	return s.viewFeedbackService
}

func (s *container) getDatabaseConnection() *sql.DB {
	if s.db != nil {
		return s.db
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
		panic(errors.New("Error while initializing connection to database"))
	}
	if err = connection.Ping(); err != nil {
		panic(errors.New("Error while pinging connection to database"))
	}
	s.db = connection
	return s.db
}

var instance Container

func GetInstance() Container {
	if instance == nil {
		New(config.GetInstance())
	}
	return instance
}
