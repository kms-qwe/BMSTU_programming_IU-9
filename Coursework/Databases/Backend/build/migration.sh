#!/bin/sh
set -eu

: "${DB_HOST:=postgres}"
: "${DB_PORT:=5432}"
: "${DB_USER:=coffee}"
: "${DB_PASSWORD:=coffee}"
: "${DB_NAME:=coffee}"
: "${DB_SSL_MODE:=disable}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

echo "Running database migrations..."
goose -dir /migrations postgres "${DATABASE_URL}" up
