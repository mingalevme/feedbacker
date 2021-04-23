package envvarbag

import (
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

type EnvVarBag interface {
	Get(key string, fallback string) string
	Require(key string) string
	GetInt(key string, fallback int) (int, error)
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

func (s *Bag) Require(key string) string {
	val, ok := s.storage[key]
	if ok && strutils.IsNonEmptyString(val) {
		return val
	} else {
		panic(errors.Errorf("Missing %s env var", key))
	}
}

func (s *Bag) GetInt(key string, fallback int) (int, error) {
	val, ok := s.storage[key]
	if !ok {
		return fallback, nil
	}
	if i, err := strconv.ParseInt(val, 10, 0); err != nil {
		return fallback, err
	} else {
		return int(i), nil
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
