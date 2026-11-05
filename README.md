# Mercado Rural

Proyecto de página web realizado para la cátedra de Programación Web de la FCEx UNICEN.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Richter, Matias

## Descripción

- ~~TP1: Esta etapa del proyecto consiste en hacer el primer servidor web local, donde escuche en el puerto 8080, y sirva un HTML estatico.~~

- ~~TP2: Esta etapa del proyecto consiste en definir una capa de persistencia mediante PostgreSQL. Incluye un pipeline de automatización para la generación de código con sqlc, montado de
una base de datos con Docker, Tests y un script para automatizar el setup.~~

- ~~TP3: Esta etapa del proyecto consiste en la construcción de una API REST en Go para exponer las entidades del dominio mediante handlers HTTP con operaciones CRUD completas (POST, GET, PUT, DELETE). Conecta la lógica de negocio con la capa de datos de sqlc, incluye manejo y validaciones JSON, estados de respuesta HTTP estándar (201, 200, 204, 404), despliegue containerizado con Docker y una batería de pruebas de integración reproducibles mediante Hurl.~~

- TP4: Interfaz de usuario en HTML, CSS y JavaScript que consume la API REST. Permite crear usuarios, animales y publicaciones, listar publicaciones activas y eliminar registros.

## Estructura del proyecto

```text
Mercadorural/
├── docker-compose.yml       // Configuración del contenedor
│
├── db/
│   ├── queries/                 // Archivos SQL para las operaciones CRUD
│   ├── schema/                  // Archivos SQL de generación de tablas
│   └── sqlc/                    // Archivos generados por sqlc y tests
├── static/
│   ├── index.html               // Página principal del frontend
│   ├── app.js                   // Lógica de la interfaz y llamadas a la API
│   └── styles.css               // Estilos de la interfaz
├── go.mod / go.sum              // Dependencias de Go
├── main.go                      // Punto de entrada de la aplicación
├── Makefile                     // Automatización de base de datos y tests
├── handlePublicaciones.go       // Handler de publicaciones
├── handleUsuario.go             // Handler de usuarios
├── handleVentas.go              // Handler de ventas
├── handleAnimales.go            // Handler de animales
└── checkPoint.hurl              // Pruebas de integración de la API
```

## Requisitos

1. [Go Compiler](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/get-docker/)
3. [Docker Compose](https://docs.docker.com/compose/)
4. [sqlc](https://sqlc.dev/)
5. [Make](https://www.gnu.org/software/make/)
6. [Hurl](https://hurl.dev/) (para las pruebas de integración)
7. [curl](https://curl.se/)

## Instalación

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/TpWeb2026/Mercadorural.git
   ```

2. Ingresar a la carpeta del proyecto:

   ```bash
   cd Mercadorural
   ```

## Cómo ejecutar la API y los tests

- `make test`: limpia el contenedor, genera el código con sqlc, levanta PostgreSQL, carga el esquema y ejecuta las pruebas unitarias de `db/sqlc`.
- `make test-api`: realiza el mismo setup, levanta el servidor en el puerto 8080, ejecuta `checkPoint.hurl` y luego detiene el servidor y la base de datos.

Ambos comandos eliminan el volumen de Docker, por lo que borran los datos cargados previamente.

## Cómo ejecutar la interfaz web (TP4)

1. Levantar PostgreSQL:

   ```bash
   docker compose up -d
   ```

2. Esperar a que la base de datos esté lista. La primera vez (o después de un `make test`) hay que cargar el esquema:

   ```bash
   cat db/schema/*.sql | docker compose exec -T database psql -U postgres -d tp2
   ```

   Si la base ya tiene las tablas creadas, este paso no es necesario.

3. Ejecutar el servidor Go desde la raíz del proyecto:

   ```bash
   go run .
   ```

4. Abrir en el navegador:

   ```text
   http://localhost:8080
   ```
