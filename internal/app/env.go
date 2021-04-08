package app

import (
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/emailer"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
	"os"
	"os/user"
	"strconv"
)

type Env interface {
	AppEnv() string
	Debug() bool
	// HTTP
	MaxPostRequestBodyLength() uint
	LogRequests() bool
	//
	FeedbackRepository() repository.Feedback
	//
	Logger() log.Logger
	//
	NotifierEmailFrom() string
	NotifierEmailTo() string
	NotifierEmailSubjectTemplate() string
	Notifier() notifier.FeedbackLeftNotifier
	//
	MailSmtpHost() string
	MailSmtpPort() uint16
	MailSmtpUsername() *string
	MailSmtpPassword() *string
	EmailSender() emailer.EmailSender
	//
	DBHost() string
	DBPort() uint16
	DBUser() string
	DBPass() string
	DBName() string
	DatabaseConnection() *sql.DB
	//
	Sentry() *sentry.Hub
	Rollbar() *rollbar.Client
}

type Container struct {
	envVarBag  envvarbag.EnvVarBag
	logger     log.Logger
	repository repository.Feedback
	notifier   notifier.FeedbackLeftNotifier
	emailer    emailer.EmailSender
	db         *sql.DB
	sentry     *sentry.Hub
	rollbar    *rollbar.Client
}

func NewEnv(e envvarbag.EnvVarBag) *Container {
	return &Container{
		envVarBag: e,
	}
}

func (s *Container) getEnvVar(key, fallback string) string {
	return s.envVarBag.Get(key, fallback)
}

func (s *Container) AppEnv() string {
	return s.envVarBag.Get("APP_ENV", "production")
}

func (s *Container) Debug() bool {
	val, err := strconv.ParseBool(s.envVarBag.Get("DEBUG", "0"))
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing DEBUG envVarBag-var"))
	}
	return val
}

func (s *Container) MaxPostRequestBodyLength() uint {
	val, err := strconv.ParseUint(s.envVarBag.Get("MAX_POST_REQUEST_BODY_LENGTH", ""), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAX_POST_REQUEST_BODY_LENGTH envVarBag-var"))
	}
	return uint(val)
}

func (s *Container) LogRequests() bool {
	v, err := strconv.ParseBool(s.envVarBag.Get("HTTP_LOG_REQUESTS", "0"))
	if err != nil {
		panic(errors.Wrap(err, "env: parsing HTTP_LOG_REQUESTS to bool"))
	}
	return v
}

func (s *Container) FeedbackRepository() repository.Feedback {
	if s.repository != nil {
		return s.repository
	}
	driver := s.envVarBag.Get("PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		conn := s.DatabaseConnection()
		// https://github.com/go-pg/pg
		s.repository = repository.NewDatabaseFeedbackRepository(conn, s.Logger())
	} else if driver == "array" {
		s.repository = repository.NewArrayFeedbackRepository(s.Logger())
	} else {
		panic(errors.Errorf("Unsupported persistence driver: %s", driver))
	}
	return s.repository
}

func (s *Container) MailSmtpHost() string {
	return s.envVarBag.Get("MAIL_SMTP_HOST", "127.0.0.1")
}

func (s *Container) MailSmtpPort() uint16 {
	val, err := strconv.ParseUint(s.envVarBag.Get("MAIL_SMTP_PORT", "25"), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAIL_SMTP_PORT envVarBag-var"))
	}
	if val > 65535 {
		panic(errors.Wrapf(err, "Value of MAIL_SMTP_PORT envVarBag-var is too big: %d", val))
	}
	return uint16(val)
}

func (s *Container) MailSmtpUsername() *string {
	u := s.envVarBag.Get("MAIL_SMTP_USERNAME", "")
	if util.IsEmptyString(u) {
		return nil
	}
	return &u
}

func (s *Container) MailSmtpPassword() *string {
	p := s.envVarBag.Get("MAIL_SMTP_PASSWORD", "")
	if util.IsEmptyString(p) {
		return nil
	}
	return &p
}

func (s *Container) NotifierEmailFrom() string {
	def := func() string {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s@%s", u.Username, h)
	}
	from := s.envVarBag.Get("NOTIFIER_EMAIL_FROM", "")
	if util.IsEmptyString(from) {
		return def()
	}
	return from
}

func (s *Container) NotifierEmailTo() string {
	to := s.envVarBag.Get("NOTIFIER_EMAIL_TO", "")
	if util.IsEmptyString(to) {
		panic(errors.New("NOTIFIER_EMAIL_TO is empty"))
	}
	return to
}

func (s *Container) EmailSender() emailer.EmailSender {
	if s.emailer != nil {
		return s.emailer
	}
	driver := s.envVarBag.Get("EMAILER_DRIVER", "smtp")
	if driver == "smtp" {
		s.emailer = emailer.NewSmtpEmailSender(s.MailSmtpHost(), s.MailSmtpPort(), s.MailSmtpUsername(), s.MailSmtpPassword(), s.Logger())
	} else if driver == "array" {
		s.emailer = emailer.NewArrayEmailSender(s.Logger())
	} else {
		panic(errors.Errorf("Unsupported emailer driver: %s", driver))
	}
	return s.emailer
}

func (s *Container) NotifierEmailSubjectTemplate() string {
	return s.envVarBag.Get("NOTIFIER_EMAIL_SUBJECT_TEMPLATE", "Feedback %{InstallationID}s")
}

func (s *Container) Notifier() notifier.FeedbackLeftNotifier {
	if s.notifier != nil {
		return s.notifier
	}
	driver := s.envVarBag.Get("NOTIFIER_DRIVER", "email")
	if driver == "email" {
		s.notifier = notifier.NewEmailFeedbackLeftNotifier(s.EmailSender(), s.NotifierEmailFrom(), s.NotifierEmailTo(), s.NotifierEmailSubjectTemplate(), s.Logger())
	} else if driver == "array" {
		s.notifier = notifier.NewArrayFeedbackLeftNotifier(s.Logger())
	} else {
		panic(errors.Errorf("Unsupported notifier driver: %s", driver))
	}
	return s.notifier
}

func (s *Container) Sentry() *sentry.Hub {
	if s.sentry != nil {
		return s.sentry
	}
	dsn := s.envVarBag.Get("SENTRY_DSN", "")
	if util.IsEmptyString(dsn) {
		panic("SENTRY_DSN-envvar is empty")
	}
	s.sentry = s.newSentryHub(sentry.ClientOptions{
		Dsn: dsn,
		Debug: s.Debug(),
		Environment: s.AppEnv(),
	})
	s.sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetFingerprint([]string{"{{ default }}", "{{ message }}", "{{ error.type }}", "{{ error.value }}"})
	})
	return s.sentry
}

func (s *Container) newSentryHub(opts sentry.ClientOptions) *sentry.Hub {
	client, err := sentry.NewClient(opts)
	if err != nil {
		panic(err)
	}
	return sentry.NewHub(client, sentry.NewScope())
}

func (s *Container) Rollbar() *rollbar.Client {
	if s.rollbar != nil {
		return s.rollbar
	}
	token := s.envVarBag.Get("ROLLBAR_TOKEN", "")
	if util.IsEmptyString(token) {
		panic("ROLLBAR_TOKEN-envvar is empty")
	}
	s.rollbar = rollbar.New(token, s.AppEnv(), "", "", "")
	s.rollbar.SetFingerprint(true)
	return s.rollbar
}
