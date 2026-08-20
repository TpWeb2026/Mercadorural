# Mi proyecto

Este es un proyecto realizado para la materia de Programacion Web.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Ritcher, Matias

## Descripción

El proyecto consiste en hacer el primer servidor web local, donde escuche en el puerto 8080.

## Estructura del proyecto
El proyecto por ahora tiene una estructura bastante simple.

```text
mi_servidor_web/
├── main.go
└── static/
    └── index.html
```

## Cómo ejecutar
Antes de poder ejecutar el proyecto, se tiene que clonar, para poder tener este proyecto localmente en la computadora

1. Clonar el repositorio.
   Para eso es necesario primero abrir una terminal y ejecutar el comando:
   git clone https://github.com/TpWeb2026/Mercadorural.git
   
3. Ubicarnos en la carpeta del proyecto.
  Para eso tenemos dos formas de hacerlo, la primera es hacerla manuelamente,
  buscando donde esta esa carpeta con el nombre de "Mercadorural", la segunda es
  abrir una terminal y ejecutar el comando cd Mercadorural.

4. Una vez ya ubicados en la carpeta, tenemos que ejecutar el siguiente comando: go run main.go
   Si todo funciona bien, en la terminal, tendria que aparecer un mensaje que diga     "Servidor escuchando en http://localhost:8080"

5. Si el mensaje anterior aparecio en la terminal, estamos en condicones de abrir      el navegador, el de tu preferencia y en la barra de busqueda, escribir lo        siguiente:
   http://localhost:8080
   Eso mostrara la pagina web que hemos creado.
   
6. Para terminar de cerrar, tenemos que detener el servidor, lo cual se hace,         volviendo a la terminal que estuvimos usando, donde nos aparece el mensaje de
  "Servidor escuchando en http://localhost:8080", y hacer la combinacion de teclas de ctrl + c.
   
