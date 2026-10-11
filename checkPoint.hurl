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
[Captures] #usamos el Captures para poder guardar el id que nos devuelve el json, sino tendriamos que borrar y crear cada vez que ejecutemos el hurl
usuario_id: jsonpath "$.id"  #guardamos el id en una variable
[Asserts]
jsonpath "$.id" exists
jsonpath "$.nombre" == "Juanito"

# LISTAR USUARIOS (GET /usuario/)
GET http://localhost:8080/usuario/
HTTP 200

# OBTENER USUARIO POR ID (GET /usuario/{{usuario_id}})
GET http://localhost:8080/usuario/{{usuario_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{usuario_id}}
jsonpath "$.nombre" == "Juanito"

# ACTUALIZAR USUARIO POR ID (PUT /usuario/{{usuario_id}})

PUT http://localhost:8080/usuario/{{usuario_id}}
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

# ELIMINAR USUARIO POR ID (DELETE /usuario/{{usuario_id}})
DELETE http://localhost:8080/usuario/{{usuario_id}}
HTTP 204

# VERIFICAR BORRADO (GET /usuario/1 debe dar 404)
GET http://localhost:8080/usuario/{{usuario_id}}
HTTP 404


# ------ Chequeamos cada endpoint de los animales ------


