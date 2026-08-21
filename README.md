# Mercado Rural

Proyecto de pagina WEB realizado para la catedra de Programación Web de la FCex UNICEN.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Richter, Matias

## Descripción

Esta etapa del proyecto consiste en hacer el primer servidor web local, donde escuche en el puerto 8080, y sirva un HTML estatico.

## Estructura del proyecto

```text
Mercadorural/
├── main.go
└── static/
    └── index.html
```

## Requisitos

1. [Go Compiler](https://go.dev/dl/)

## Cómo ejecutar

1. Clonar el repositorio.
   En una terminal (Bash, CMD, Powershell) ejecutar el comando:

   ```
   git clone https://github.com/TpWeb2026/Mercadorural.git
   ```

2. Ubicarnos en la carpeta del proyecto.

   ```
   cd Mercadorural
   ```

3. Una vez ya ubicados en la carpeta, tenemos que ejecutar el siguiente comando:

   ```
   go run main.go
   ```

   Si todo funciona bien, en la terminal, mostrara el mensaje:  "Servidor escuchando en <http://localhost:8080>"

4. Abrir el Navegador, en la direccion local:

   ```
   http://localhost:8080
   ```

5. Para terminar la ejecucion del servidor, usar la combinacion de teclas CTRL + C en la terminal.
