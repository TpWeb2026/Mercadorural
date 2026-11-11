# Mercado Rural

Proyecto de página web realizado para la cátedra de Programación Web de la FCEx UNICEN.

## Integrantes

- Allende, Ignacio Matias
- Picchioni, Patricio
- Richter, Matias

## Descripción

Mercado Rural es una aplicación web para la compra y venta de animales. Los usuarios pueden registrarse, cargar animales y publicarlos con un precio. Otro usuario puede comprar una publicación; la operación queda registrada en la tabla `Venta` y la publicación pasa a mostrarse como vendida, sin volver a estar disponible.

Esta entrega corresponde al **TP5**, que reemplaza la arquitectura de API JSON + JavaScript por una aplicación renderizada en el servidor con [templ](https://templ.guide/). Las vistas se generan en Go y los formularios funcionan con recargas completas de página y el patrón Post/Redirect/Get.

La entrega se realiza en la rama `tp5`.

## Arquitectura

- **Go + `net/http`**: servidor y ruteo.
- **templ**: componentes HTML renderizados en el servidor.
- **PostgreSQL + sqlc**: persistencia y acceso a datos.
- **Pico.css**: estilos base, vendorizado en `static/pico.min.css`.
- **`static/styles.css`**: estilos propios del proyecto.
- **Sin JavaScript**: la aplicación no contiene archivos `.js`.

La entidad `Venta` funciona como historial inmutable: se registra al comprar una publicación y no se expone ninguna acción para editarla ni eliminarla. La publicación comprada queda inactiva y se muestra como vendida.

## Estructura del proyecto

```text
Mercadorural/
├── docker-compose.yml
├── db/
│   ├── queries/
│   ├── schema/
│   └── sqlc/
├── views/
│   ├── layout.templ
│   ├── index.templ
│   ├── entity_form.templ
│   ├── entity_list.templ
│   ├── usuario.templ
│   ├── animal.templ
│   ├── venta.templ
│   ├── view_models.go
│   └── *_templ.go
├── static/
│   ├── pico.min.css
│   └── styles.css
├── main.go
├── server.go
├── handleUsuario.go
├── handleAnimales.go
├── handlePublicaciones.go
├── helpers.go
├── run-app.sh
├── Makefile
├── go.mod / go.sum
└── README.md
```

## Requisitos

1. [Go](https://go.dev/dl/) 1.27 o superior.
2. [Docker](https://docs.docker.com/get-docker/) con Docker Compose v2.
3. [Make](https://www.gnu.org/software/make/).
4. [sqlc](https://sqlc.dev/) disponible en el `PATH` o en `$(go env GOPATH)/bin/sqlc`.

`templ` se instala como dependencia de herramienta del módulo y se ejecuta con `go tool templ`; no hace falta instalarlo globalmente.

## Cómo ejecutar

Desde la raíz del repositorio:

```bash
make run-app
```

También se puede ejecutar directamente:

```bash
bash run-app.sh
```

El comando o script realiza automáticamente:

1. Borra contenedores y volúmenes previos.
2. Ejecuta `sqlc generate`.
3. Ejecuta `go tool templ generate`.
4. Compila la aplicación con `go build`.
5. Levanta PostgreSQL con Docker Compose.
6. Espera a que la base esté disponible.
7. Carga el esquema de la base.
8. Ejecuta la aplicación en `http://localhost:8080`.
9. Al detener la aplicación, borra los contenedores y volúmenes.

La aplicación queda disponible en:

```text
http://localhost:8080
```

Para detenerla, usar `Ctrl+C`; el script se encarga de la limpieza final.

### Ejecución manual

```bash
docker compose down -v
go tool templ generate
go run .
```

En ese caso, la base y el esquema deben prepararse por separado con Docker Compose.

## Rutas

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/` | Página principal con formularios y listados |
| POST | `/usuario` | Crea un usuario |
| POST | `/animales` | Crea un animal |
| POST | `/publicacion` | Crea una publicación |
| POST | `/publicacion/{id}/comprar` | Registra una venta y marca la publicación como vendida |
| POST | `/usuario/{id}/eliminar` | Da de baja un usuario |
| POST | `/animales/{id}/eliminar` | Elimina un animal |
| POST | `/publicacion/{id}/eliminar` | Elimina una publicación |

Los formularios usan `POST` y redirigen con `303 See Other` a la sección correspondiente de la página principal.

## Orden recomendado de uso

1. Crear un usuario vendedor.
2. Crear un animal seleccionando el usuario como dueño.
3. Crear una publicación seleccionando el animal y el vendedor.
4. Crear un segundo usuario para que actúe como comprador.
5. En la publicación activa, seleccionar el comprador y presionar **Comprar**.
6. Verificar que la publicación quede marcada como **Vendida** y que la operación aparezca en **Ventas realizadas**.
7. Probar el borrado de registros no vendidos.

## Pruebas y verificación

Para ejecutar las pruebas unitarias con una base limpia:

```bash
make test
```

Para regenerar el código de templ y sqlc:

```bash
make generate
```

Para formatear el código Go y los archivos templ:

```bash
make fmt
```

Para bajar contenedores y volúmenes:

```bash
make db-down
```

## Notas de implementación

- Los archivos `views/*_templ.go` son generados por templ y no deben editarse a mano.
- Después de modificar un archivo `.templ`, ejecutar `go tool templ generate` o `make generate`.
- Los precios se envían como texto porque sqlc representa las columnas `NUMERIC` como `string` en Go.
- `DeleteUsuario` realiza una baja lógica (`activo = false`), mientras que los animales y las publicaciones se eliminan físicamente.
- La compra se ejecuta en una transacción: se marca la publicación como inactiva y se inserta la venta. Si algo falla, ambas operaciones se revierten.
- El comprador debe ser un usuario activo distinto del vendedor.
- Una publicación vendida deja de mostrar el formulario de compra y no puede volver a comprarse.
- Las ventas son un historial inmutable: no tienen acciones de edición ni borrado.
- Si un animal o una publicación tiene ventas asociadas, la base rechaza el borrado y la aplicación muestra un mensaje de error.
- Pico.css está vendorizado para que la aplicación funcione sin conexión a internet.
- No se incluyen archivos JavaScript en la entrega.

## Entrega

La entrega se sube al repositorio en la rama `tp5`.
