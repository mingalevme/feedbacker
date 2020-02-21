# Feedbacker
Example of a simple HTTP API application written in Go (router, database + migrations, json request/response, task queues)

# Run PostgreSQL
```
docker run -d --restart=always --name feedbacker-postgres postgres
# or
#docker run -d --restart=always -p "54321:5432" --name feedbacker-postgres postgres
# or
#docker run -d --restart=always -p "${FEEDBACKER_DB_PORT}:5432" --name "${FEEDBACKER_DB_HOST}" postgres
```

# Remove PostgreSQL
```
docker rm --force feedbacker-postgresql
```

# Migations

## Create migration
```
docker run -rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations create -ext sql -dir /migrations -seq create_feedback_table
```

## Run migrations
```
docker run --rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations -database "postgres://postgres@localhost:5432/postgres?sslmode=disable" up
# Or via linking
#docker run --rm -v "${PWD}/db/migrations:/migrations" --link="feedbacker-postgres:postgres" migrate/migrate -path=/migrations -database "postgres://postgres@postgres/postgres?sslmode=disable" up
```

## Rollback all migrations
```
docker run --rm -v "${PWD}/db/migrations:/migrations" --network host migrate/migrate -path=/migrations -database "postgres://postgres@localhost:54322/postgres?sslmode=disable" down -all
```

# TODO
- https://echo.labstack.com/guide
- https://github.com/google/go-cloud/tree/master/wire
