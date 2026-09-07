-- name: CreateUsuario :one
INSERT INTO Usuario (nombre, apellido, contacto)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsuario :one
SELECT id, nombre, apellido, contacto 
FROM Usuario 
WHERE id = $1 LIMIT 1;