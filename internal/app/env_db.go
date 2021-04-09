package app

import (
	"database/sql"
	"github.com/mingalevme/feedbacker/pkg/util"
	"github.com/pkg/errors"
	"strconv"
)

func (s *Container) DBDriver() string {
	return s.EnvVarBag.Get("DB_DRIVER", "postgres")
}

func (s *Container) DBHost() string {
	return s.EnvVarBag.Get("DB_HOST", "127.0.0.1")
}

func (s *Container) DBPort() uint16 {
	val, err := strconv.ParseUint(s.EnvVarBag.Get("DB_PORT", "5432"), 10, 0)
	if err != nil {
		panic(errors.Wrap(err, "Error while parsing MAIL_SMTP_PORT env-var"))
	}
	if val > 65535 {
		panic(errors.Wrapf(err, "Value of MAIL_SMTP_PORT env-var is too big: %d", val))
	}
	return uint16(val)
}

func (s *Container) DBUser() string {
	return s.EnvVarBag.Get("DB_USER", "postgres")
}

func (s *Container) DBPass() string {
	return s.EnvVarBag.Get("DB_PASS", "postgres")
}

func (s *Container) DBName() string {
	return s.EnvVarBag.Get("DB_NAME", "postgres")
}

func (s *Container) DatabaseConnection() *sql.DB {
	if s.db != nil {
		return s.db
	}
	params := map[string]interface{}{
		"Driver":   s.DBDriver(),
		"Host":     s.DBHost(),
		"Port":     strconv.Itoa(int(s.DBPort())),
		"User":     s.DBUser(),
		"Pass":     s.DBPass(),
		"Database": s.DBName(),
	}
	dataSourceName := util.Sprintf("%{Driver}s://%{User}s:%{Pass}s@%{Host}s:%{Port}s/%{Database}s?sslmode=disable", params)
	connection, err := sql.Open(s.DBDriver(), dataSourceName)
	if err != nil {
		panic(errors.New("Error while initializing connection to database"))
	}
	if err = connection.Ping(); err != nil {
		panic(errors.Wrap(err, "pinging connection to database"))
	}
	s.db = connection
	return s.db
}
