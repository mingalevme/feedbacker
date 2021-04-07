package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/_ddd/infrastructure/env"
	"github.com/mingalevme/feedbacker/internal/app/di"
	"github.com/mingalevme/feedbacker/internal/config"
	"os"
)

func main() {
	cfg := config.New(env.New())
	container := di.New(cfg)
	conn := container.GetDatabaseConnection()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("usage: %s dir", os.Args[0]))
	}

	path := os.Args[1]

	driver, err := postgres.WithInstance(conn, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		cfg.GetDBName(),
		driver,
	)

	if err != nil {
		panic(err)
	}

	if err := migrator.Up(); err != nil {
		panic(err)
	}
}
