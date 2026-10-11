# ------ Chequeamos cada endpoint de los usuarios ------

# CREAR USUARIO (POST /usuario/)
POST http://localhost:8080/usuario/
Content-Type: application/json
{
  "nombre": "Juanito",
  "apellido": "Tortuguita",
  "contacto": "OlA12@email.com"
}
HTTP 201
[Asserts]
jsonpath "$.id" exists
jsonpath "$.nombre" == "Juanito"

# LISTAR USUARIOS (GET /usuario/)
GET http://localhost:8080/usuario/
HTTP 200

# OBTENER USUARIO POR ID (GET /usuario/1)
GET http://localhost:8080/usuario/1
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.nombre" == "Juanito"

# ACTUALIZAR USUARIO POR ID (PUT /usuario/1)

PUT http://localhost:8080/usuario/1
Content-Type: application/json
{
  "nombre": "HOla tubi",
  "apellido": "hola1 123",
  "contacto": "ggg.actualiggggzado@email.com"
}

HTTP 200
[Asserts]
jsonpath "$.nombre" == "HOla tubi"
jsonpath "$.apellido" == "hola1 123"

# ELIMINAR USUARIO POR ID (DELETE /usuario/1)
DELETE http://localhost:8080/usuario/1
HTTP 204

# VERIFICAR BORRADO (GET /usuario/1 debe dar 404)
GET http://localhost:8080/usuario/1
HTTP 404


# ------ Chequeamos cada endpoint de los animales ------
