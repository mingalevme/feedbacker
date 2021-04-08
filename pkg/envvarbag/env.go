package envvarbag

import (
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type EnvVarBag interface {
	Get(key string, fallback string) string
}

type Bag struct {
	storage map[string]string
}

func (s *Bag) Get(key string, fallback string) string {
	val, ok := s.storage[key]
	if ok {
		return val
	} else {
		return fallback
	}
}

func New() *Bag {
	// Make a copy of environment
	m := map[string]string{}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		m[variable[0]] = variable[1]
	}
	return &Bag{
		storage: m,
	}
}

func NewWithValues(values map[string]string) *Bag {
	b := New()
	for k, v := range values {
		b.storage[k] = v
	}
	return b
}
