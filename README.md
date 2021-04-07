# Feedbacker
Example of a simple HTTP API application written in Go (router, database + migrations, json request/response)

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

## Run app


# TODO
- https://echo.labstack.com/guide
- https://github.com/google/go-cloud/tree/master/wire
