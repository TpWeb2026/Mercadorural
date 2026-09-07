-- name: CreateVenta :one
INSERT INTO Venta (id_publicacion, id_vendedor, id_comprador)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVenta :one
SELECT id, id_publicacion, id_vendedor, id_comprador, fecha 
FROM Venta 
WHERE id = $1 LIMIT 1;

-- name: ListVentasByComprador :many
SELECT id, id_publicacion, id_vendedor, id_comprador, fecha 
FROM Venta 
WHERE id_comprador = $1 
ORDER BY fecha DESC;

-- name: ListVentasConDetalle :many
SELECT v.id AS venta_id, v.fecha, 
       u_comp.nombre AS comprador_nombre, 
       u_vend.nombre AS vendedor_nombre, 
       p.precio
FROM Venta v
JOIN Usuario u_comp ON v.id_comprador = u_comp.id
JOIN Usuario u_vend ON v.id_vendedor = u_vend.id
JOIN Publicacion p ON v.id_publicacion = p.id;