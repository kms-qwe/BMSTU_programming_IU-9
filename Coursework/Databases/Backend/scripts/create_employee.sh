#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/build/.env"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "env file not found: ${ENV_FILE}" >&2
  exit 1
fi

if [[ $# -ne 4 ]]; then
  echo "usage: $0 <full_name> <phone> <login> <password>" >&2
  exit 1
fi

FULL_NAME="$1"
PHONE="$2"
LOGIN="$3"
PASSWORD="$4"

set -a
# shellcheck disable=SC1090
source "${ENV_FILE}"
set +a

docker compose \
  --env-file "${ENV_FILE}" \
  -f "${ROOT_DIR}/build/docker-compose.yml" \
  exec -T postgres \
  psql \
    -v ON_ERROR_STOP=1 \
    -U "${DB_USER}" \
    -d "${DB_NAME}" \
    -c "INSERT INTO employees (full_name, phone, login, password_hash) VALUES ('${FULL_NAME//\'/\'\'}', '${PHONE//\'/\'\'}', '${LOGIN//\'/\'\'}', crypt('${PASSWORD//\'/\'\'}', gen_salt('bf')));"

echo "employee created: ${LOGIN}"
