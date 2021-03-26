package db

import (
	"database/sql"
	"github.com/mingalevme/feedbacker/infrastructure/env"
	"github.com/mingalevme/feedbacker/infrastructure/utils"
	"github.com/pkg/errors"
)

func NewDatabaseConnection() (*sql.DB, error) {
	params := map[string]interface{}{
		"Host":     env.GetEnvValue("FEEDBACKER_DB_HOST", "127.0.0.1"),
		"Port":     env.GetEnvValue("FEEDBACKER_DB_PORT", "5432"),
		"User":     env.GetEnvValue("FEEDBACKER_DB_USER", "postgres"),
		"Pass":     env.GetEnvValue("FEEDBACKER_DB_PASSWORD", "postgres"),
		"Database": env.GetEnvValue("FEEDBACKER_DB_NAME", "postgres"),
	}
	dataSourceName := utils.Sprintf("postgres://%{User}s:%{Pass}s@%{Host}s:%{Port}s/%{Database}s?sslmode=disable", params)
	connection, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "Error while initializing connection to database")
	}
	if err = connection.Ping(); err != nil {
		return nil, errors.Wrap(err, "Error while pinging connection to database")
	}
	return connection, nil
}
