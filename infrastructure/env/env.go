package env

import (
	_ "github.com/lib/pq"
	"os"
	"strings"
)

var instance *env

type Env interface {
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

// Initializer!
func GetEnv() Env {
	if instance == nil {
		// Make a copy of environment
		m := map[string]string{}
		for _, element := range os.Environ() {
			variable := strings.Split(element, "=")
			m[variable[0]] = variable[1]
		}
		instance = &env{
			environment: m,
		}
	}
	return instance
}

// Quick Accessor
func GetEnvValue(key string, fallback string) string {
	return GetEnv().Get(key, fallback)
}
