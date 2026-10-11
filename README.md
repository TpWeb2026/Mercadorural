# Mercado Rural

Proyecto de pagina WEB realizado para la catedra de Programación Web de la FCex UNICEN.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Richter, Matias

## Descripción

- ~~TP1: Esta etapa del proyecto consiste en hacer el primer servidor web local, donde escuche en el puerto 8080, y sirva un HTML estatico.~~

- ~~TP2: Esta etapa del proyecto define una capa de persistencia mediante PostgreSQL. Incluye un pipeline de automatización para la generación de código con sqlc, montado de
una base de datos con Docker, Tests y un script para automatizar el setup.~~

- TP3: Esta etapa del proyecto consiste en la construcción de una API REST en Go para exponer las entidades del dominio mediante handlers HTTP con operaciones CRUD completas (POST, GET, PUT, DELETE). Conecta la lógica de negocio con la capa de datos de sqlc, incluye manejo y validaciones JSON, estados de respuesta HTTP estándar (201, 200, 204, 404), despliegue containerizado con Docker y una batería de pruebas de integración reproducibles mediante Hurl.

## Estructura del proyecto

```text
Mercadorural/
├── docker-compose.yml       // Configuracion del contenedor
│
├── db/
│   ├── queries/                 // Archivos SQL para las operaciones CRUD
│   ├── schema/                  // Archivos SQL de generacion de tablas
│   └── sqlc/                    // Archivos generados por sqlc y tests
├── static/
│   └── index.html               // HTML de entrada
├── go.mod / go.sum              // Dependencias de go
├── main.go                      // Punto de entrada de la aplicacion
├── Makefile                     // Script de deploy y test
├── handlePublicaciones.go       // Handle para la entidad de publicaciones
├── handleUsuario.go             // Handle para la entidad de usuario
├── handleVentas.go              // Handle para la entidad de ventas
├── handleAnimales.go            // Handle para la entidad de animales
└── checkPoint.hurl              // Archivo donde se prueba cada endpoint de cada entidad 

```

## Requisitos

1. [Go Compiler](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/get-docker/)
3. Sqlc generate
4. Docker-Compose
5. Make
6. Hurl
7. Curl

## Instalacion

1. Clonar el repositorio.
   En una terminal (Bash, CMD, Powershell) ejecutar el comando:

   ```
   git clone https://github.com/TpWeb2026/Mercadorural.git
   ```
   Si ya lo tenes clonado de la entrega pasada, solo queda pararse en esa carpeta y hacer git pull para que traiga todo lo nuevo que se hizo 


## Cómo ejecutar

1. Ubicarnos en la carpeta del proyecto.

   ```
   cd Mercadorural
   ```

2. Una vez ya ubicados en la carpeta, tenemos que ejecutar el siguiente comando:

   ```
   make test-api
   ```
   Con ese comando, se ejecuta el archivo MakeFile, donde en el mismo se encuentra toda la automatizacion
   para limpiar el contendor, generar el codigo go con sqlc generate, levantar el docker, esperar a que se 
   levante la base de datos, una vez que el docker esta levantado, se ejecuta el main para que se levante el servidor y se ponga a escuchar en el puerto 8080. Cuando se levanto el servidor, se ejecuta el hurl, el cual chequea que anden correctamente cada endpoint de cada entidad, una vez realizado eso, detiene el servidor y da de baja la base de datos.