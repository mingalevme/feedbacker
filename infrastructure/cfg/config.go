package cfg

import (
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/pkg/errors"
	"strconv"
)

type Config interface {
	//GetEnvBag() env.EnvVarBag
	GetEnvVar(key string, fallback string) string
	IsDebug() bool
	GetMaxPostRequestBodyLength() uint
	GetHTTPHost() string
	GetHTTPPort() string
}

type config struct {
	envVarBag env.EnvVarBag
}

func New(e env.EnvVarBag) Config {
	return &config{
		envVarBag: e,
	}
}

func (s *config) GetHTTPHost() string {
	return s.envVarBag.Get("HTTP_HOST", "localhost")
}

func (s *config) GetHTTPPort() string {
	return s.envVarBag.Get("HTTP_PORT", "8080")
}

//func (s *cfg) GetEnvBag() env.EnvVarBag {
//	return s.envVarBag
//}

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
