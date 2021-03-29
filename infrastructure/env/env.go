package env

import (
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type EnvVarBag interface {
	Get(key string, fallback string) string
}

type env struct {
	environment map[string]string
}

func (s *env) Get(key string, fallback string) string {
	val, ok := s.environment[key]
	if ok {
		return val
	} else {
		return fallback
	}
}

func New() EnvVarBag {
	// Make a copy of environment
	m := map[string]string{}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		m[variable[0]] = variable[1]
	}
	return &env{
		environment: m,
	}
}

//
//var instance EnvVarBag
//
//func GetDefault() EnvVarBag {
//	if instance == nil {
//		instance = New()
//	}
//	return instance
//}
//
// Quick Accessor
//func GetEnvValue(key string, fallback string) string {
//	return GetDefault().Get(key, fallback)
//}
