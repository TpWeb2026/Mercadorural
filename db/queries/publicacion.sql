-- name: CreatePublicacion :one
INSERT INTO Publicacion (id_animal, precio, id_vendedor)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPublicacion :one
SELECT id, id_animal, precio, id_vendedor 
FROM Publicacion 
WHERE id = $1 LIMIT 1;

-- name: ListPublicacionesActivas :many
SELECT id, id_animal, precio, id_vendedor 
FROM Publicacion 
ORDER BY id DESC;