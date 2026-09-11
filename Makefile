COMPOSE_FILE=compose/docker-compose.yml
DB_CONTAINER=database
DB_USER=postgres
DB_NAME=tp2
SCHEMA_DIR=db/schema

# .PHONY se usa para declarar que las palabras son comandos
.PHONY: all pre test post clean

all: pre test post

#tareas previas
pre:
	@echo "-> Generando codigo Go con sqlc generate..."
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
	
	@echo "-> Compilando el proyecto..."
	go build ./...
	
	@echo "-> Limpiando contenedores previos ..."
	docker compose -f $(COMPOSE_FILE) down -v
	
	@echo "-> Iniciando el contenedor de la base de datos PosgreSQL..."
	docker compose -f $(COMPOSE_FILE) up -d
	
	@echo "-> Levantando la base de datos..."
	@while ! docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 2; \
	done
	@echo "-> Base de datos: UP"
	
	@echo "-> Insertando esquemas en la base de datos..."
	@cat $(SCHEMA_DIR)/*.sql | docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) > /dev/null
	@echo "-> Tareas previas finalizadas!"
	@echo ""

# tests
test:
	@echo "Ejecutando TEST"
	go test -v -count=1 ./db/sqlc/...
	@echo ""

# tareas posteriores
post: clean

# limpieza final
clean:
	@echo "Limpiando entorno"
	@echo "-> Deteniendo contenedores y limpiando volumenes"
	docker compose -f $(COMPOSE_FILE) down -v
	@echo "-> Entorno limpio"