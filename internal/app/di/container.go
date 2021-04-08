package di

import (
	"database/sql"
	"github.com/getsentry/sentry-go"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/emailer"
	"github.com/mingalevme/feedbacker/internal/app/service/log"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/internal/config"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
	"strconv"
)

type Container interface {
	GetFeedbackRepository() repository.Feedback
	GetLogger() log.Logger
	//GetLeaveFeedbackService() LeaveFeedback
	//GetViewFeedbackService() service.ViewFeedbackService
	GetFeedbackLeftNotifier() notifier.FeedbackLeftNotifier
	GetEmailSender() emailer.EmailSender
	GetDatabaseConnection() *sql.DB
}

type container struct {
	config     config.Config
	logger     log.Logger
	repository repository.Feedback
	notifier   notifier.FeedbackLeftNotifier
	emailer    emailer.EmailSender
	db         *sql.DB
	sentry     *sentry.Hub
	rollbar    *rollbar.Client
	//leaveFeedbackService service.LeaveFeedbackService
	//viewFeedbackService  service.ViewFeedbackService
}

func New(config config.Config) Container {
	instance = &container{
		config: config,
	}
	return instance
}

func (s *container) GetSentry() *sentry.Hub {
	if s.sentry != nil {
		return s.sentry
	}
	dsn := s.config.GetEnvVar("SENTRY_DSN", "")
	if util.IsEmptyString(dsn) {
		panic("SENTRY_DSN-envvar is empty")
	}
	s.sentry = s.newSentryHub(sentry.ClientOptions{
		Dsn: dsn,
		Debug: s.config.IsDebug(),
		Environment: s.config.GetAppEnv(),
	})
	s.sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetFingerprint([]string{"{{ default }}", "{{ message }}", "{{ error.type }}", "{{ error.value }}"})
	})
	return s.sentry
}

func (s *container) newSentryHub(opts sentry.ClientOptions) *sentry.Hub {
	client, err := sentry.NewClient(opts)
	if err != nil {
		panic(err)
	}
	return sentry.NewHub(client, sentry.NewScope())
}

func (s *container) GetRollbar() *rollbar.Client {
	if s.rollbar != nil {
		return s.rollbar
	}
	token := s.config.GetEnvVar("ROLLBAR_TOKEN", "")
	if util.IsEmptyString(token) {
		panic("ROLLBAR_TOKEN-envvar is empty")
	}
	s.rollbar = rollbar.New(token, s.config.GetAppEnv(), "", "", "")
	s.rollbar.SetFingerprint(true)
	return s.rollbar
}

func (s *container) newRollbarClient(opts sentry.ClientOptions) *sentry.Hub {
	client, err := sentry.NewClient(opts)
	if err != nil {
		panic(err)
	}
	return sentry.NewHub(client, sentry.NewScope())
}

func (s *container) GetFeedbackRepository() repository.Feedback {
	if s.repository != nil {
		return s.repository
	}
	driver := s.config.GetEnvVar("PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		conn := s.GetDatabaseConnection()
		// https://github.com/go-pg/pg
		s.repository = repository.NewDatabaseFeedbackRepository(conn, s.GetLogger())
	} else if driver == "array" {
		s.repository = repository.NewArrayFeedbackRepository(s.GetLogger())
	} else {
		panic(errors.Errorf("Unsupported persistence driver: %s", driver))
	}
	return s.repository
}

func (s *container) GetDatabaseConnection() *sql.DB {
	if s.db != nil {
		return s.db
	}
	params := map[string]interface{}{
		"Host":     s.config.GetDBHost(),
		"Port":     strconv.Itoa(int(s.config.GetDBPort())),
		"User":     s.config.GetDBUser(),
		"Pass":     s.config.GetDBPass(),
		"Database": s.config.GetDBName(),
	}
	dataSourceName := util.Sprintf("postgres://%{User}s:%{Pass}s@%{Host}s:%{Port}s/%{Database}s?sslmode=disable", params)
	connection, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(errors.New("Error while initializing connection to database"))
	}
	if err = connection.Ping(); err != nil {
		panic(errors.Wrap(err, "pinging connection to database"))
	}
	s.db = connection
	return s.db
}

func (s *container) GetFeedbackLeftNotifier() notifier.FeedbackLeftNotifier {
	if s.notifier != nil {
		return s.notifier
	}
	driver := s.config.GetEnvVar("NOTIFIER_DRIVER", "email")
	if driver == "email" {
		s.notifier = notifier.NewEmailFeedbackLeftNotifier(s.GetEmailSender(), s.config.GetNotifierEmailFrom(), s.config.GetNotifierEmailTo(), s.config.GetNotifierEmailSubjectTemplate(), s.GetLogger())
	} else if driver == "array" {
		s.notifier = notifier.NewArrayFeedbackLeftNotifier(s.GetLogger())
	} else {
		panic(errors.Errorf("Unsupported notifier driver: %s", driver))
	}
	return s.notifier
}

func (s *container) GetEmailSender() emailer.EmailSender {
	if s.emailer != nil {
		return s.emailer
	}
	driver := s.config.GetEnvVar("EMAILER_DRIVER", "smtp")
	if driver == "smtp" {
		s.emailer = emailer.NewSmtpEmailSender(s.config.GetMailSmtpHost(), s.config.GetMailSmtpPort(), s.config.GetMailSmtpUsername(), s.config.GetMailSmtpPassword(), s.GetLogger())
	} else if driver == "array" {
		s.emailer = emailer.NewArrayEmailSender(s.GetLogger())
	} else {
		panic(errors.Errorf("Unsupported emailer driver: %s", driver))
	}
	return s.emailer
}

//func (s *container) GetLeaveFeedbackService() service.LeaveFeedbackService {
//	if s.leaveFeedbackService == nil {
//		s.leaveFeedbackService = service.NewLeaveFeedbackService(s.GetFeedbackRepository(), s.GetLogger())
//	}
//	return s.leaveFeedbackService
//}
//
//func (s *container) GetViewFeedbackService() service.ViewFeedbackService {
//	if s.viewFeedbackService == nil {
//		s.viewFeedbackService = service.NewViewFeedbackService(s.GetFeedbackRepository(), s.GetLogger())
//	}
//	return s.viewFeedbackService
//}

var instance Container

func GetInstance() Container {
	if instance == nil {
		New(config.GetInstance())
	}
	return instance
}
