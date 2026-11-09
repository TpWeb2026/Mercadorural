SHELL := /bin/bash

COMPOSE_FILE=docker-compose.yml
DB_CONTAINER=database
DB_USER=postgres
DB_NAME=tp2
SCHEMA_DIR=db/schema

GO ?= go
TEMPL ?= $(GO) tool templ
SQLC ?= $(shell command -v sqlc 2>/dev/null || printf '%s/bin/sqlc' "$$($(GO) env GOPATH)")

.PHONY: generate fmt test run-app run dev db-down clean

generate:
	$(TEMPL) generate
	$(SQLC) generate

fmt:
	gofmt -w .
	$(TEMPL) fmt .

test:
	@bash -eu -c '\
		cleanup() { docker compose -f $(COMPOSE_FILE) down -v >/dev/null 2>&1 || true; }; \
		trap cleanup EXIT; \
		docker compose -f $(COMPOSE_FILE) down -v; \
		$(TEMPL) generate; \
		$(SQLC) generate; \
		docker compose -f $(COMPOSE_FILE) up -d; \
		while ! docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT 1" >/dev/null 2>&1; do sleep 1; done; \
		cat $(SCHEMA_DIR)/*.sql | docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) >/dev/null; \
		$(GO) test -v -count=1 ./...; \
	'

run-app:
	bash ./run-app.sh

run: run-app

dev: run-app

db-down:
	docker compose -f $(COMPOSE_FILE) down -v

clean: db-down
