package config

import (
	"fmt"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/env"
	util2 "github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"os"
	"os/user"
	"strconv"
)

type Config interface {
	//GetEnvBag() env.EnvVarBag
	GetEnvVar(key string, fallback string) string
	IsDebug() bool
	GetMaxPostRequestBodyLength() uint
	GetMailSmtpHost() string
	GetMailSmtpPort() uint16
	GetMailSmtpUsername() *string
	GetMailSmtpPassword() *string
	GetNotifierEmailFrom() string
	GetNotifierEmailTo() string
	GetNotifierEmailSubjectTemplate() string
	GetDBHost() string
	GetDBPort() uint16
	GetDBUser() string
	GetDBPass() string
	GetDBName() string
}

type config struct {
	Config
	envVarBag env.EnvVarBag
}

var instance Config

func New(e env.EnvVarBag) Config {
	instance = &config{
		envVarBag: e,
	}
	return instance
}

func (s *config) GetEnvVar(key string, fallback string) string {
	return s.envVarBag.Get(key, fallback)
}

func (s *config) IsDebug() bool {
	val, err := strconv.ParseBool(s.envVarBag.Get("DEBUG", "0"))
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing DEBUG envVarBag-var"))
	}
	return val
}

func (s *config) GetMaxPostRequestBodyLength() uint {
	val, err := strconv.ParseUint(s.envVarBag.Get("MAX_POST_REQUEST_BODY_LENGTH", ""), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAX_POST_REQUEST_BODY_LENGTH envVarBag-var"))
	}
	return uint(val)
}

func (s *config) GetMailSmtpHost() string {
	return s.envVarBag.Get("MAIL_SMTP_HOST", "127.0.0.1")
}

func (s *config) GetMailSmtpPort() uint16 {
	val, err := strconv.ParseUint(s.envVarBag.Get("MAIL_SMTP_PORT", "25"), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAIL_SMTP_PORT envVarBag-var"))
	}
	if val > 65535 {
		panic(errors.Wrapf(err, "Value of MAIL_SMTP_PORT envVarBag-var is too big: %d", val))
	}
	return uint16(val)
}

func (s *config) GetMailSmtpUsername() *string {
	u := s.envVarBag.Get("MAIL_SMTP_USERNAME", "")
	if util2.IsEmptyString(u) {
		return nil
	}
	return &u
}

func (s *config) GetMailSmtpPassword() *string {
	p := s.envVarBag.Get("MAIL_SMTP_PASSWORD", "")
	if util2.IsEmptyString(p) {
		return nil
	}
	return &p
}

func (s *config) GetNotifierEmailFrom() string {
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
	if util2.IsEmptyString(from) {
		return def()
	}
	return from
}

func (s *config) GetNotifierEmailTo() string {
	to := s.envVarBag.Get("NOTIFIER_EMAIL_TO", "")
	if util2.IsEmptyString(to) {
		panic("NOTIFIER_EMAIL_TO is empty")
	}
	return to
}

func (s *config) GetNotifierEmailSubjectTemplate() string {
	return s.envVarBag.Get("NOTIFIER_EMAIL_SUBJECT_TEMPLATE", "Feedback %{InstallationID}s")
}

func (s *config) GetDBHost() string {
	return s.GetEnvVar("DB_HOST", "127.0.0.1")
}

func (s *config) GetDBPort() uint16 {
	val, err := strconv.ParseUint(s.envVarBag.Get("DB_PORT", "5432"), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAIL_SMTP_PORT envVarBag-var"))
	}
	if val > 65535 {
		panic(errors.Wrapf(err, "Value of MAIL_SMTP_PORT envVarBag-var is too big: %d", val))
	}
	return uint16(val)
}

func (s *config) GetDBUser() string {
	return s.GetEnvVar("DB_USER", "postgres")
}

func (s *config) GetDBPass() string {
	return s.GetEnvVar("DB_PASSWORD", "postgres")
}

func (s *config) GetDBName() string {
	return s.GetEnvVar("DB_NAME", "postgres")
}

func GetInstance() Config {
	if instance == nil {
		return New(env.New())
	}
	return instance
}
