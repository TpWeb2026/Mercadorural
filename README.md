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

- ~~TP4: Interfaz de usuario en HTML, CSS y JavaScript que consume la API REST. Permite crear usuarios, animales y publicaciones, listar publicaciones activas y eliminar registros.~~

- TP5: Esta etapa reemplaza la arquitectura de API JSON + JavaScript por una aplicación renderizada en el servidor con templ. Las vistas se generan en Go y los formularios funcionan con recargas completas de página y el patrón Post/Redirect/Get.

## Flujo de la aplicación

Mercado Rural es una aplicación web para la compra y venta de animales.

1. El usuario operador se selecciona desde la barra superior. El selector funciona como placeholder del login; si no existe ningún usuario, se puede crear uno nuevo.
2. El operador publica un animal. La publicación crea el animal asociado: el animal no existe por separado.
3. El animal queda con el operador como dueño y con el mismo precio de la publicación.
4. El operador puede comprar publicaciones activas de otros usuarios, pero no las propias.
5. Al comprar, se registra una venta y la publicación pasa a estar finalizada.
6. El operador puede editar o eliminar sus publicaciones mientras estén activas. Las publicaciones finalizadas no se modifican ni se eliminan.
7. El operador puede eliminar su cuenta. Sus publicaciones activas se eliminan junto con sus animales; las publicaciones finalizadas se conservan para mantener la consistencia de las ventas.

## Arquitectura

- **Go + `net/http`**: servidor, ruteo y manejo de formularios.
- **templ**: componentes HTML renderizados en el servidor.
- **PostgreSQL + sqlc**: persistencia y acceso a datos.
- **Docker Compose**: ejecución de PostgreSQL.
- **Pico.css**: estilos base.
- **`static/styles.css`**: estilos propios del proyecto.
- **Sin JavaScript**: la aplicación no contiene archivos `.js`.

La tabla `Venta` funciona como historial inmutable: se registra al comprar una publicación y no tiene acciones de edición ni borrado. El historial no se muestra como panel en la interfaz por el momento.

## Estructura del proyecto

```text
Mercadorural/
├── docker-compose.yml       // Configuración del contenedor
│
├── db/
│   ├── queries/                 // Archivos SQL para las operaciones CRUD
│   ├── schema/                  // Archivos SQL de generación de tablas
│   └── sqlc/                    // Archivos generados por sqlc y tests
├── views/
│   ├── layout.templ             // Estructura base de la página
│   ├── index.templ              // Página principal
│   ├── entity_form.templ        // Formularios de publicación y edición
│   ├── entity_list.templ        // Listados de publicaciones propias y ajenas
│   ├── usuario.templ            // Alta de usuarios
│   ├── view_models.go           // Modelos usados por las vistas
│   └── *_templ.go               // Código generado por templ
├── static/
│   ├── pico.min.css             // Estilos base
│   └── styles.css               // Estilos propios
├── go.mod / go.sum              // Dependencias de Go
├── main.go                      // Punto de entrada de la aplicación
├── server.go                    // Rutas, operador y renderizado
├── handleUsuario.go             // Handlers de usuarios
├── handlePublicaciones.go       // Handlers de publicaciones y compras
├── helpers.go                   // Validaciones y utilidades
├── Makefile                     // Automatización de ejecución y tests
└── README.md
```

## Requisitos

1. [Go Compiler](https://go.dev/dl/)
2. [Docker](https://docs.docker.com/get-docker/)
3. [Docker Compose](https://docs.docker.com/compose/)
4. [sqlc](https://sqlc.dev/) disponible en el `PATH`.
5. [templ](https://templ.guide/) disponible en el `PATH`:

   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```

6. [Make](https://www.gnu.org/software/make/)

El Makefile invoca `sqlc generate` y `templ generate`, por lo que ambos comandos deben estar disponibles en el `PATH`.

## Instalación

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/TpWeb2026/Mercadorural.git
   ```

2. Ingresar a la carpeta del proyecto:

   ```bash
   cd Mercadorural
   ```

## Cómo ejecutar

Desde la raíz del proyecto:

```bash
make run-api
```

El comando realiza automáticamente:

1. Borra contenedores y volúmenes previos.
2. Ejecuta `sqlc generate`.
3. Ejecuta `templ generate`.
4. Compila la aplicación con `go build`.
5. Levanta PostgreSQL con Docker Compose.
6. Espera a que la base de datos esté disponible.
7. Carga el esquema de la base.
8. Ejecuta la aplicación en `http://localhost:8080`.
9. Al detener la aplicación, borra los contenedores y volúmenes.

La aplicación queda disponible en:

```text
http://localhost:8080
```

Para detenerla, usar `Ctrl+C`; el script se encarga de la limpieza final.

## Cómo probar el flujo

1. Agregar un usuario con el botón **Agregar usuario**.
2. Publicar un animal desde **Mis publicaciones**.
3. Agregar un segundo usuario y seleccionarlo como operador.
4. Comprar la publicación del primer usuario desde **Publicaciones de otros usuarios**.
5. Verificar que la publicación comprada desaparezca de los listados (queda finalizada en la base).
6. Probar la eliminación de una publicación activa y verificar que también se elimine el animal asociado.
7. Probar la eliminación de una cuenta y verificar que sus publicaciones activas desaparezcan y sus publicaciones finalizadas se conserven en la base.

## Pruebas

Para ejecutar las pruebas unitarias de la capa de datos con una base limpia:

```bash
make test
```

## Notas de implementación

- Los archivos `views/*_templ.go` son generados por templ y no deben editarse a mano.
- Después de modificar un archivo `.templ`, ejecutar `templ generate` o `make run-api`.
- Los precios se envían como texto porque sqlc representa las columnas `NUMERIC` como `string` en Go.
- El usuario operador se guarda en una cookie del servidor.
- La creación de una publicación y su animal se ejecuta en una transacción.
- La compra se ejecuta en una transacción: se marca la publicación como finalizada y se inserta la venta.
- Las publicaciones finalizadas se conservan en la base pero no se muestran en la interfaz por el momento.
- La baja de cuenta se ejecuta en una transacción: elimina las publicaciones activas y sus animales, conserva las finalizadas y marca al usuario como inactivo.
- Pico.css está vendorizado para que la aplicación funcione sin conexión a internet.
- No se incluyen archivos JavaScript en la entrega.
