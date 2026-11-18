#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

COMPOSE_FILE="docker-compose.yml"
DB_CONTAINER="database"
DB_USER="postgres"
DB_NAME="tp2"
SCHEMA_DIR="db/schema"

TMP_DIR=""

cleanup() {
  if [ -n "${TMP_DIR:-}" ]; then
    rm -rf "$TMP_DIR"
  fi

  echo "--- Limpiando contenedores y volúmenes ---"
  docker compose -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true
}

trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

if ! command -v docker >/dev/null 2>&1; then
  echo "Error: Docker no está instalado o no está en el PATH."
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "Error: Docker Compose v2 no está disponible."
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "Error: Go no está instalado o no está en el PATH."
  exit 1
fi

if command -v sqlc >/dev/null 2>&1; then
  SQLC="sqlc"
else
  SQLC="$(go env GOPATH)/bin/sqlc"
fi

if [ ! -x "$SQLC" ]; then
  echo "Error: sqlc no está instalado. Instalalo con: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
  exit 1
fi

echo "--- Borrando contenedores y volúmenes previos ---"
docker compose -f "$COMPOSE_FILE" down -v >/dev/null 2>&1 || true

echo "--- Generando código con sqlc ---"
"$SQLC" generate

echo "--- Generando código con templ ---"
go tool templ generate

echo "--- Compilando aplicación ---"
TMP_DIR="$(mktemp -d)"
APP_BIN="$TMP_DIR/mercadorural"
go build -o "$APP_BIN" .

echo "--- Levantando PostgreSQL ---"
docker compose -f "$COMPOSE_FILE" up -d

echo "--- Esperando a que PostgreSQL esté listo ---"
DB_READY=0
for _ in $(seq 1 60); do
  if docker compose -f "$COMPOSE_FILE" exec -T "$DB_CONTAINER" \
    psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" >/dev/null 2>&1; then
    DB_READY=1
    break
  fi
  sleep 1
done

if [ "$DB_READY" -ne 1 ]; then
  echo "Error: PostgreSQL no respondió a tiempo."
  exit 1
fi

echo "--- Cargando esquema de la base de datos ---"
cat "$SCHEMA_DIR"/*.sql | docker compose -f "$COMPOSE_FILE" exec -T "$DB_CONTAINER" \
  psql -U "$DB_USER" -d "$DB_NAME" >/dev/null

echo "--- Ejecutando aplicación ---"
echo "Abrí http://localhost:8080 en el navegador."
"$APP_BIN"
