COMPOSE_FILE=compose/docker-compose.yml
DB_CONTAINER=database
DB_USER=postgres
DB_NAME=tp2
SCHEMA_DIR=db/schema

.PHONY: test
 # eso se usa para evitar que si hay carpetas o archivos llamado test generen conflictos, diciendo que make es un comando y no un archivo

 #se crea la etiqueta test, para que cuando queramos ejecutar make test, se ejecute solo la etiqueta que se llama test, o sea solo esta
test:
# se usa el @ para que no imprima toda la linea, sino lo que esta dentro de echo

# Como primer paso, tenemos que limpiar todo el contenedor
	@echo "--- Limpiando contenedores y volumenes previos ---"


# Como segundo paso, tenemos que generar el codigo go con el comando sqlc generate
	@echo "--- Generando codigo Go con SQLC ---"

	sqlc generate

# Como terecer paso, ahora levantamos el contenedor con el comando docker compose up -d (se usa el -d para que no bloquee la consola)
	@echo "--- Levantando contenedor de PostgreSQL ---"

	docker compose -f $(COMPOSE_FILE) up -d

# Como cuarto paso, tenemos que esperar a que se levante el postgresql, por eso se usa este comando
	@echo "-> Levantando la base de datos..."
	@while ! docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 2; \
	done
	@echo "-> Base de datos: UP"

# Como quinto paso, tenemos que insertar todo el esquema que nos hizo en la base de datos 
	@echo "-> Insertando esquemas en la base de datos..."
	@cat $(SCHEMA_DIR)/*.sql | docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) > /dev/null

# Como sexto paso, ejecutamos el test para ver si esta bien todo lo que hicimos
	@echo "--- Ejecutando pruebas unitarias ---"

	go test -v -count=1 ./db/sqlc/...

# Como septimo paso, limpiamos todo el entorno
	@echo "--- Limpiando el entorno de Docker ---"

	docker compose -f $(COMPOSE_FILE) down -v
