COMPOSE_FILE=docker-compose.yml
DB_CONTAINER=database
DB_USER=postgres
DB_NAME=tp2
SCHEMA_DIR=db/schema

.PHONY: test  test-api
 # eso se usa para evitar que si hay carpetas o archivos llamado test generen conflictos, diciendo que make es un comando y no un archivo

 #se crea la etiqueta test, para que cuando queramos ejecutar make test, se ejecute solo la etiqueta que se llama test, o sea solo esta
test:
# se usa el @ para que no imprima toda la linea, sino lo que esta dentro de echo

# Como primer paso, tenemos que limpiar todo el contenedor
	@echo "--- Limpiando contenedores y volumenes previos ---"

	docker compose -f $(COMPOSE_FILE) down -v
	
# Como segundo paso, tenemos que generar el codigo go con el comando sqlc generate
	@echo "--- Generando codigo Go con SQLC ---"

	sqlc generate

# Como terecer paso, ahora levantamos el contenedor con el comando docker compose up -d (se usa el -d para que no bloquee la consola)
	@echo "--- Levantando contenedor de PostgreSQL ---"

	docker compose -f $(COMPOSE_FILE) up -d

# Como cuarto paso, tenemos que esperar a que se levante el postgresql, por eso se usa este comando
	@echo "-> Levantando la base de datos..."
	
	@while ! docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT 1" > /dev/null 2>&1; do \
        sleep 1; \
    done
	@echo "-> Base de datos: UP"

# Como quinto paso, tenemos que insertar todo el esquema que nos hizo en la base de datos 
	@echo "-> Insertando esquemas en la base de datos..."
	@cat $(SCHEMA_DIR)/*.sql | docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) > /dev/null

# Como sexto paso, ejecutamos el test para ver si esta bien todo lo que hicimos
	@echo "--- Ejecutando pruebas unitarias ---"

	go test -v -count=1 ./db/sqlc/...

# Como septimo paso, damos de baja el docker

	@echo "--- Limpiando el entorno de Docker ---"
	docker compose -f $(COMPOSE_FILE) down -v


# -------------- PRUEBAS DE LA INTEGRACION DE LA API -------------

test-api:

# Usamos las barras invertidas para que no abra una terminal nueva cada vez que se lee una linea, si no las usariamos, perderiamos las variables donde se guarda el id de las cosas
# Con timeout le damos un tiempo de 10 ciclos para que se levante el servidor, si eso no pasa se termina de ejecutar todo

	@echo "--- Limpiando entorno previo ---"
	docker compose -f $(COMPOSE_FILE) down -v

	@echo "--- Generando código Go con SQLC ---"
	sqlc generate

	@echo "--- Levantando PostgreSQL ---"
	docker compose -f $(COMPOSE_FILE) up -d

	@echo "-> Esperando a que PostgreSQL esté listo..."
	@while ! docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT 1" > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "-> Base de datos: UP"

	@echo "-> Cargando esquemas SQL..."
	@cat $(SCHEMA_DIR)/*.sql | docker compose -f $(COMPOSE_FILE) exec -T $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) > /dev/null

	@echo "--- Levantando API y ejecutando pruebas Hurl ---"
	@go run . & ID_SERVIDOR=$$!; \
	echo "-> Esperando respuesta de http://localhost:8080..."; \
	TIMEOUT=10; \
	while [ $$TIMEOUT -gt 0 ] && ! curl -s http://localhost:8080/usuario/ > /dev/null 2>&1; do \
		sleep 0.5; \
		TIMEOUT=$$((TIMEOUT - 1)); \
	done; \
	if [ $$TIMEOUT -eq 0 ]; then \
		echo "El servidor Go no respondió a tiempo."; \
		docker compose -f $(COMPOSE_FILE) down -v; \
		exit 1; \
	fi; \
	echo "--- Ejecutando pruebas de integración Hurl ---"; \
	hurl --test checkPoint.hurl; \
	HURL_STATUS=$$?; \
	echo "--- Deteniendo servidor Go ---"; \
	kill $$ID_SERVIDOR 2>/dev/null || true; \
	echo "--- Limpiando el entorno de Docker ---" \
	docker compose -f $(COMPOSE_FILE) down -v; \
	if [ $$HURL_STATUS -ne 0 ]; then \
		echo "Fallaron las pruebas de Hurl."; \
		exit $$HURL_STATUS; \
	fi; \
	echo "Pruebas de integración finalizadas con éxito"
