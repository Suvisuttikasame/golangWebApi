#!/bin/sh
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "postgres://postgres:postgres@golangwebapi-postgres-1:5432/simple_bank?sslmode=disable" -verbose up

echo "start the app"
# take all parameter and run in this case will be '/app/myapp' 
exec "$@"