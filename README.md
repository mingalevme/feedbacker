# Feedbacker
Example of a simple HTTP API application written in Go (router, database + migrations, json request/response)

## Stack

- [Project structure](https://github.com/golang-standards/project-layout)
- Separating http and app logics
- Environment based Config + Dependency Injection: ./internal/app/env.go
- HTTP Middleware [Echo](https://echo.labstack.com/):
  - HTTP Server
  - Routing
  - Data binding
- [Database migrations](https://github.com/golang-migrate/migrate)
- (Multichannel) Contextable (data / error / request) logger / interface
  - Null
  - Stdout (Logrus based)
  - [Logrus](https://github.com/sirupsen/logrus)
  - [Rollbar](https://rollbar.com/)
  - [Sentry](https://sentry.io/)
  - Stack (multichannel)
  - Abstract Logger to implement custom logger
  - StdLog Adapter (TODO)
- Multi driver data repository
  - Null driver
  - Array driver
  - Database driver
  - Redis driver
- Email sending
  - Null
  - Array
  - SMTP
- Testing
  - Testing only build (go-) tags: [testing.go](https://github.com/mingalevme/feedbacker/blob/master/internal/app/model/testing.go)
  - [Mocking database](https://github.com/mingalevme/feedbacker/blob/master/internal/app/repository/database_test.go) via [DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
  - [Mocking redis](https://github.com/mingalevme/feedbacker/blob/master/internal/app/repository/redis_test.go) via [elliotchance/redismock](https://github.com/elliotchance/redismock)

# Examples of running PostgreSQL via Docker
```
docker run -d --restart=always -e POSTGRES_PASSWORD=postgres --name feedbacker-postgres postgres
# or
docker run -d --restart=always -p "54321:5432" -e POSTGRES_PASSWORD=postgres --name feedbacker-postgres postgres
# or
docker run -d --restart=always -p "${FEEDBACKER_DB_PORT}:5432" -e POSTGRES_PASSWORD=postgres --name "${FEEDBACKER_DB_HOST}" postgres
```

# Examples of checking PostgreSQL connection via Docker
docker run --rm -it --network host postgres psql -h 127.0.0.1 -p 54321 -U "postgres" postgres

# Examples of (force) deleting PostgreSQL via Docker
```
docker rm --force feedbacker-postgres
```

# Migrations (https://github.com/golang-migrate/migrate) via Docker

## Create a migration
```
docker run -rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations create -ext sql -dir /migrations -seq create_feedback_table
```

## Run migrations
```
DB_PORT=54321 go run cmd/migrate.go internal/db/migrations
# or
docker run --rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations -database "postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
# Or via linking (depreacted, use custom networks)
docker run --rm -v "${PWD}/db/migrations:/migrations" --link="feedbacker-postgres:postgres" migrate/migrate -path=/migrations -database "postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
```

## Rollback all migrations
```
docker run --rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations -database "postgres://postgres@localhost:54322/postgres?sslmode=disable" down -all
```

## Run app (examples)

### Local/Debug mode

```go build -gcflags="all=-N -l"``` build with remote debug support


#### Memory persistence driver
```
go build -gcflags="all=-N -l" && PERSISTENCE_DRIVER=array ./feedbacker
```

#### Database (PostgreSQL) persistence driver
```
go build -gcflags="all=-N -l" && PERSISTENCE_DRIVER=database DB_HOST=127.0.0.1 DB_PORT=5432 DB_USER=feedbacker DB_PASS=feedbacker DB_NAME=feedbacker ./feedbacker
```

#### Redis persistence driver
```
go build -gcflags="all=-N -l" && PERSISTENCE_DRIVER=redis REDIS_ADDR="127.0.0.1:6379" REDIS_PASS="xxx" REDIS_DB=1 ./feedbacker
```

#### Logging stdout channel
```
go build -gcflags="all=-N -l" && LOG_CHANNEL=stdout LOG_STDOUT_LEVEL=debug ./feedbacker
```

#### Logging Sentry channel
```
go build -gcflags="all=-N -l" && LOG_CHANNEL=sentry LOG_SENTRY_LEVEL=warning SENTRY_DSN="https://...ingest.sentry.io/..." ./feedbacker
```

#### Logging Rollbar channel
```
go build -gcflags="all=-N -l" && LOG_CHANNEL=rollbar LOG_ROLLBAR_LEVEL=error ROLLBAR_TOKEN="..." ./feedbacker
```

#### Logging stack channel
```
go build -gcflags="all=-N -l" && LOG_CHANNEL=stack LOG_STACK_CHANNELS="sentry,rollbar,stdout" LOG_SENTRY_LEVEL=warning LOG_ROLLBAR_LEVEL=error LOG_STDOUT_LEVEL=debug SENTRY_DSN="https://...ingest.sentry.io/..." ROLLBAR_TOKEN="..." ./feedbacker
```

#### Logging null channel
```
go build -gcflags="all=-N -l" && LOG_CHANNEL=null ./feedbacker
```

#### Notifying email channel
```
go build -gcflags="all=-N -l" && NOTIFIER_DRIVER=email NOTIFIER_EMAIL_TO=user@example.com ./feedbacker
```

#### Notifying memory channel
```
go build -gcflags="all=-N -l" && NOTIFIER_DRIVER=array ./feedbacker
```

#### Notifying null channel
```
go build -gcflags="all=-N -l" && NOTIFIER_DRIVER=null ./feedbacker
```

### Testing
```
go test -v -cover -tags testing ./...
```

### Docker (build & run)

```
docker build -t feedbacker . && docker run -e "PERSISTENCE_DRIVER=array" -e "LOG_CHANNEL=null" -e "NOTIFIER_DRIVER=null" -e "..." -it feedbacker
```
