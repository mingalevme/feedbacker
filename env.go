package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"fmt"
	"github.com/pkg/errors"
)

type Env struct {
	environ                  map[string]string
	db                       *sql.DB
	logger                   *log.Logger
	debug                    *bool
	maxPostRequestBodyLength *uint
}

func NewEnv() *Env {
	// Make a copy of environment
	m := map[string]string{}
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		m[variable[0]] = variable[1]
	}
	env := &Env{
		environ: m,
	}
	return env
}

func (e *Env) GetEnv(key string, fallback string) string {
	val, ok := e.environ[key]
	if ok && val != "" {
		e.Logger().Debugf("Environment variable \"%s\" is set: \"%s\"", key, val)
		return val
	} else {
		e.Logger().Debugf("Environment variable \"%s\" does not set, using fallback value \"%s\"", key, fallback)
		return fallback
	}
}

func (e *Env) Db() *sql.DB {
	if e.db == nil {
		e.db = e.initDb()
	}
	return e.db
}

func (e *Env) initDb() *sql.DB {
	params := map[string]interface{}{
		"Host":     e.GetEnv("FEEDBACKER_DB_HOST", "127.0.0.1"),
		"Port":     e.GetEnv("FEEDBACKER_DB_PORT", "5432"),
		"User":     e.GetEnv("FEEDBACKER_DB_USER", "postgres"),
		"Pass":     e.GetEnv("FEEDBACKER_DB_PASSWORD", "postgres"),
		"Database": e.GetEnv("FEEDBACKER_DB_NAME", "postgres"),
	}
	dataSourceName := Sprintf("postgres://%{User}s:%{Pass}s@%{Host}s:%{Port}s/%{Database}s?sslmode=disable", params)
	e.Logger().Infof("Connecting to database: %s", dataSourceName)
	connection, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		e.Logger().Fatal("Error while initializing connection to database ", err)
	}
	if err = connection.Ping(); err != nil {
		e.Logger().Fatal("Error while pinging connection to database ", err)
	} else {
		e.Logger().Info("Connection to database has been established successfully")
	}
	return connection
}

func (e *Env) Logger() *log.Logger {
	if e.logger == nil {
		e.logger = log.New()
		e.logger.SetOutput(os.Stdout)
		if level, err := log.ParseLevel(e.GetEnv("FEEDBACKER_LOG_LEVEL", "debug")); err != nil {
			fmt.Errorf("%+v", errors.Wrap(err, "Error while parsing log level"))
			e.logger.SetLevel(log.DebugLevel)
		} else {
			e.logger.SetLevel(level)
		}
	}
	return e.logger
}

func (e *Env) Debug() bool {
	if e.debug == nil {
		value, _ := strconv.ParseBool(e.GetEnv("FEEDBACKER_DEBUG", "0"))
		e.debug = &value
	}
	return *e.debug
}

func (e *Env) MaxPostRequestBodyLength() uint {
	if e.maxPostRequestBodyLength == nil {
		value64, err := strconv.ParseUint(e.GetEnv("FEEDBACKER_MAX_POST_REQUEST_BODY_LENGTH", ""), 10, 0)
		if err != nil {
			value64 = uint64(1024 * 1024)
		}
		value := uint(value64)
		e.maxPostRequestBodyLength = &value
	}
	return *e.maxPostRequestBodyLength
}
