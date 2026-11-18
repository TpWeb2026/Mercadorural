-- name: CreateUsuario :one
INSERT INTO Usuario (nombre, apellido, contacto)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsuario :one
SELECT id, nombre, apellido, contacto 
FROM Usuario 
WHERE id = $1 AND activo = true LIMIT 1;

-- name: UpdateUsuario :one
UPDATE Usuario
SET nombre = $2,
    apellido = $3,
    contacto = $4
WHERE id = $1
RETURNING id, nombre, apellido, contacto, activo;

-- name: DeleteUsuario :exec
UPDATE Usuario
SET activo = false      /*flag para borrado logico, y poder mantener los valores
                         transaccionales de Venta.*/
WHERE id = $1;

-- name: ListUsuarios :many
SELECT id, nombre, apellido, contacto, activo
FROM Usuario WHERE activo = true 
ORDER BY id;