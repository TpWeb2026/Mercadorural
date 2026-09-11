# Mercado Rural

Proyecto de pagina WEB realizado para la catedra de Programación Web de la FCex UNICEN.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Richter, Matias

## Descripción

- ~~TP1: Esta etapa del proyecto consiste en hacer el primer servidor web local, donde escuche en el puerto 8080, y sirva un HTML estatico.~~

- TP2: Esta etapa del proyecto define una capa de persistencia mediante PostgreSQL. Incluye un pipeline de automatización para la generación de código con sqlc, montado de
una base de datos con Docker, Tests y un script para automatizar el setup

## Estructura del proyecto

```text
Mercadorural/
├── compose/
│   └── docker-compose.yml       // Configuracion del contenedor
├── db/
│   ├── queries/                 // Archivos SQL para las operaciones CRUD
│   ├── schema/                  // Archivos SQL de generacion de tablas
│   └── sqlc/                    // Archivos generados por sqlc y tests
├── static/
│   └── index.html               // HTML de entrada
├── go.mod / go.sum              // Dependencias de go
├── main.go                      // Punto de entrada de la aplicacion
└── Makefile                     // Script de deploy y test
```

## Requisitos

1. [Go Compiler](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/get-docker/)
3. Sqlc generate
4. Docker-Compose
5. Make

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
   go run main.go
   ```

   Si todo funciona bien, en la terminal, mostrara el mensaje:  "Servidor escuchando en <http://localhost:8080>"

3. Abrir el Navegador, en la direccion local:

   ```
   http://localhost:8080
   ```

4. Para terminar la ejecucion del servidor, usar la combinacion de teclas CTRL + C en la terminal.

## Cómo ejecutar Tests

   1. Ubicarnos en la carpeta del proyecto.

      ```
      cd Mercadorural
      ```

   2. Ejecutar el script de deploy

      ```
      make
      ```