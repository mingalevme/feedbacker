// https://martinfowler.com/articles/injection.html#UsingAServiceLocator

package app

import (
	"context"
	"database/sql"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/mingalevme/feedbacker/internal/app/repository"
	"github.com/mingalevme/feedbacker/internal/app/service/notifier"
	"github.com/mingalevme/feedbacker/pkg/dispatcher"
	"github.com/mingalevme/feedbacker/pkg/emailer"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
	"github.com/mingalevme/feedbacker/pkg/log"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
	"github.com/slack-go/slack"
	"strconv"
)

type Env interface {
	EnvironmentID() string
	Debug() bool

	MaxPostRequestBodyLength() uint
	LogRequests() bool

	Logger() log.Logger

	FeedbackRepository() repository.Feedback

	NotifierEmailFrom() string
	NotifierEmailTo() string
	NotifierEmailSubjectTemplate() string
	Notifier() notifier.Notifier

	Dispatcher() dispatcher.Dispatcher
	TaskQueue() dispatcher.TaskQueue

	MailSmtpHost() string
	MailSmtpPort() uint16
	MailSmtpUsername() *string
	MailSmtpPassword() *string
	EmailSender() emailer.EmailSender

	DBDriver() string
	DBHost() string
	DBPort() uint16
	DBUser() string
	DBPass() string
	DBName() string
	DatabaseConnection() *sql.DB

	Sentry() *sentry.Hub
	Rollbar() *rollbar.Client

	RedisAddr() string
	RedisPass() string
	RedisDB() uint
	Redis() *redis.Client

	Slack() *slack.Client

	Build()
	Close()
}

type Container struct {
	EnvVarBag  envvarbag.EnvVarBag
	logger     log.Logger
	dispatcher dispatcher.Dispatcher
	taskQueue  dispatcher.TaskQueue
	repository repository.Feedback
	notifier   notifier.Notifier
	emailer    emailer.EmailSender
	db         *sql.DB
	sentry     *sentry.Hub
	rollbar    *rollbar.Client
	redis      *redis.Client
	slack      *slack.Client
}

func NewEnv(e envvarbag.EnvVarBag) *Container {
	return &Container{
		EnvVarBag: e,
	}
}

func (s *Container) getEnvVar(key, fallback string) string {
	return s.EnvVarBag.Get(key, fallback)
}

func (s *Container) EnvironmentID() string {
	return s.EnvVarBag.Get("APP_ENV", "production")
}

func (s *Container) Debug() bool {
	val, err := strconv.ParseBool(s.EnvVarBag.Get("DEBUG", "0"))
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing DEBUG env-var"))
	}
	return val
}

func (s *Container) MaxPostRequestBodyLength() uint {
	val, err := strconv.ParseUint(s.EnvVarBag.Get("MAX_POST_REQUEST_BODY_LENGTH", ""), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAX_POST_REQUEST_BODY_LENGTH env-var"))
	}
	return uint(val)
}

func (s *Container) LogRequests() bool {
	v, err := strconv.ParseBool(s.EnvVarBag.Get("HTTP_LOG_REQUESTS", "0"))
	if err != nil {
		panic(errors.Wrap(err, "env: parsing HTTP_LOG_REQUESTS to bool"))
	}
	return v
}

func (s *Container) TaskQueue() dispatcher.TaskQueue {
	if s.taskQueue != nil {
		return s.taskQueue
	}
	s.taskQueue = dispatcher.NewTaskQueue(s.Dispatcher())
	return s.taskQueue
}

func (s *Container) Dispatcher() dispatcher.Dispatcher {
	if s.dispatcher != nil {
		return s.dispatcher
	}
	driver := s.EnvVarBag.Get("DISPATCHER_DRIVER", "chan")
	if driver == "chan" {
		queueMaxSize, _ := s.EnvVarBag.GetInt("DISPATCHER_CHAN_MAX_QUEUE_SIZE", 100)
		workersCount, _ := s.EnvVarBag.GetInt("DISPATCHER_CHAN_WORKER_COUNT", 1)
		s.dispatcher = dispatcher.NewChanDriver(s.Logger(), queueMaxSize, workersCount)
	} else if driver == "sync" {
		s.dispatcher = dispatcher.NewSyncDriver(s.Logger())
	} else if driver == "go" {
		maxRunningProcessCount, _ := s.EnvVarBag.GetInt("DISPATCHER_GO_MAX_RUNNING_PROCESS_COUNT", 1)
		s.dispatcher = dispatcher.NewGoDriver(s.Logger(), maxRunningProcessCount)
	} else if driver == "array" {
		s.dispatcher = dispatcher.NewArrayDriver()
	} else if driver == "null" {
		s.dispatcher = dispatcher.NewNullDriver(s.Logger())
	} else {
		panic(errors.Errorf("Unsupported dispatcher driver: %s", driver))
	}
	return s.dispatcher
}

func (s *Container) FeedbackRepository() repository.Feedback {
	if s.repository != nil {
		return s.repository
	}
	driver := s.EnvVarBag.Get("PERSISTENCE_DRIVER", "database")
	if driver == "database" {
		conn := s.DatabaseConnection()
		// https://github.com/go-pg/pg
		s.repository = repository.NewDatabaseFeedbackRepository(conn, s.Logger())
	} else if driver == "redis" {
		s.repository = repository.NewRedisFeedbackRepository(s.Redis(), context.Background())
	} else if driver == "array" {
		s.repository = repository.NewArrayFeedbackRepository(s.Logger())
	} else {
		panic(errors.Errorf("Unsupported persistence driver: %s", driver))
	}
	return s.repository
}

func (s *Container) MailSmtpHost() string {
	return s.EnvVarBag.Get("MAIL_SMTP_HOST", "127.0.0.1")
}

func (s *Container) MailSmtpPort() uint16 {
	val, err := strconv.ParseUint(s.EnvVarBag.Get("MAIL_SMTP_PORT", "25"), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAIL_SMTP_PORT env-var"))
	}
	if val > 65535 {
		panic(errors.Wrapf(err, "Value of MAIL_SMTP_PORT env-var is too big: %d", val))
	}
	return uint16(val)
}

func (s *Container) MailSmtpUsername() *string {
	u := s.EnvVarBag.Get("MAIL_SMTP_USERNAME", "")
	if util.IsEmptyString(u) {
		return nil
	}
	return &u
}

func (s *Container) MailSmtpPassword() *string {
	p := s.EnvVarBag.Get("MAIL_SMTP_PASSWORD", "")
	if util.IsEmptyString(p) {
		return nil
	}
	return &p
}

func (s *Container) EmailSender() emailer.EmailSender {
	if s.emailer != nil {
		return s.emailer
	}
	driver := s.EnvVarBag.Get("EMAILER_DRIVER", "smtp")
	if driver == "smtp" {
		s.emailer = emailer.NewSmtpEmailSender(s.MailSmtpHost(), s.MailSmtpPort(), s.MailSmtpUsername(), s.MailSmtpPassword(), s.Logger())
	} else if driver == "sendmail" {
		cmd := s.EnvVarBag.Require("EMAILER_SENDMAIL_CMD")
		s.emailer = emailer.NewSendmailEmailSender(cmd, s.Logger())
	} else if driver == "array" {
		s.emailer = emailer.NewArrayEmailSender()
	} else if driver == "null" {
		s.emailer = emailer.NewNullEmailSender()
	} else {
		panic(errors.Errorf("Unsupported emailer driver: %s", driver))
	}
	return s.emailer
}

func (s *Container) Slack() *slack.Client {
	if s.slack != nil {
		return s.slack
	}
	t := s.EnvVarBag.Get("SLACK_TOKEN", "")
	if strutils.IsEmptyString(t) {
		panic("Missing SLACK_TOKEN env var")
	}
	s.slack = slack.New(t)
	return s.slack
}

func (s *Container) Sentry() *sentry.Hub {
	if s.sentry != nil {
		return s.sentry
	}
	dsn := s.EnvVarBag.Get("SENTRY_DSN", "")
	if util.IsEmptyString(dsn) {
		panic("Missing SENTRY_DSN env var")
	}
	s.sentry = s.newSentryHub(sentry.ClientOptions{
		Dsn:         dsn,
		Debug:       s.Debug(),
		Environment: s.EnvironmentID(),
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
	token := s.EnvVarBag.Get("ROLLBAR_TOKEN", "")
	if util.IsEmptyString(token) {
		panic("ROLLBAR_TOKEN-envvar is empty")
	}
	s.rollbar = rollbar.New(token, s.EnvironmentID(), "", "", "")
	s.rollbar.SetFingerprint(true)
	return s.rollbar
}

func (s *Container) RedisAddr() string {
	return s.getEnvVar("REDIS_ADDR", "127.0.0.1:6379")
}

func (s *Container) RedisPass() string {
	return s.getEnvVar("REDIS_PASS", "")
}

func (s *Container) RedisDB() uint {
	v := s.getEnvVar("REDIS_DB", "0")
	n, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing REDIS_DB env-var"))
	}
	return uint(n)
}

func (s *Container) Redis() *redis.Client {
	if s.redis != nil {
		return s.redis
	}
	s.redis = redis.NewClient(&redis.Options{
		Addr:     s.RedisAddr(),
		Password: s.RedisPass(),
		DB:       int(s.RedisDB()),
	})
	return s.redis
}

func (s *Container) Build() {
	s.TaskQueue()
	s.FeedbackRepository()
	s.Logger()
	s.Notifier()
}

func (s *Container) Close() {
	if s.logger != nil {
		s.logger.Close()
	}
}
