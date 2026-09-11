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

-- name: UpdateAnimal :one
UPDATE Animal
SET nombre = $2,
    raza = $3,
    id_dueno = $4,
    precio = $5
WHERE id = $1
RETURNING *;

-- name: DeleteAnimal :exec
DELETE FROM Animal
WHERE id = $1;

-- name: ListAnimales :many
SELECT id, nombre, raza, id_dueno, precio 
FROM Animal 
ORDER BY id;
