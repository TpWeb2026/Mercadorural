-- name: CreateAnimal :one
INSERT INTO Animal (nombre, raza, id_dueno, precio)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAnimal :one
SELECT id, nombre, raza, id_dueno, precio 
FROM Animal 
WHERE id = $1 LIMIT 1;

-- name: ListAnimalesByDueno :many
SELECT id, nombre, raza, id_dueno, precio 
FROM Animal 
WHERE id_dueno = $1;
