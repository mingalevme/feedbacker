package envvarbag

import (
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type EnvVarBag interface {
	Get(key string, fallback string) string
	With(values map[string]string) EnvVarBag
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

func (s *Bag) With(values map[string]string) EnvVarBag {
	bag := &Bag{
		storage: map[string]string{},
	}
	for k, v := range s.storage {
		 bag.storage[k] = v
	}
	if values != nil {
		for k, v := range values {
			bag.storage[k] = v
		}
	}
	return bag
}

func New() *Bag {
	// Make a copy of environment
	storage := map[string]string{}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		storage[variable[0]] = variable[1]
	}
	return &Bag{
		storage: storage,
	}
}

func NewWithValues(values map[string]string) *Bag {
	b := New()
	for k, v := range values {
		b.storage[k] = v
	}
	return b
}
