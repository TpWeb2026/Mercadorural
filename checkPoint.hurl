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


# Antes de empezar a chequear los endpoint, tenemos que crear un usuario, porque animales tiene que tener un dueño si o si 

#Creamos el vendedor del animal
POST http://localhost:8080/usuario/
Content-Type: application/json
{
  "nombre": "Pedro",
  "apellido": "Dueño",
  "contacto": "pedro@email.com"
}

HTTP 201
[Captures] #aca capturamos el id, porque no sabemos cuantos usuarios hay, se hace para poder hacerlo aun con la base con mas usuarios y que sea mas robusto el hurl
vendedor_id: jsonpath "$.id"


# Empezamos a hacer el test de los endpoint de la entidad de animal

# CREAR ANIMAL (POST /animales/)
POST http://localhost:8080/animales/
Content-Type: application/json
{
  "nombre": "Firulais",
  "raza": "Ovejero Alemán",
  "id_dueno": {{vendedor_id}},
  "precio": "150.00"
}
HTTP 201
[Captures] # una vez que creamos el animal, capturamos su id, porque puede haber mas animales en la base de datos
animal_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.nombre" == "Firulais"
jsonpath "$.id_dueno" == {{vendedor_id}}

# LISTAR ANIMALES (GET /animales/)
GET http://localhost:8080/animales/
HTTP 200

# OBTENER ANIMAL POR ID (GET /animales/{{animal_id}})
GET http://localhost:8080/animales/{{animal_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{animal_id}}
jsonpath "$.nombre" == "Firulais"

# ACTUALIZAR ANIMAL POR ID (PUT /animales/{{animal_id}})
PUT http://localhost:8080/animales/{{animal_id}}
Content-Type: application/json
{
  "nombre": "Firulais Editado",
  "raza": "Ovejero Alemán Puro",
  "id_dueno": {{vendedor_id}},
  "precio": "200.00"
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Firulais Editado"
jsonpath "$.precio" == "200.00"

# ELIMINAR ANIMAL POR ID (DELETE /animales/{{animal_id}})
DELETE http://localhost:8080/animales/{{animal_id}}
HTTP 204

# VERIFICAR BORRADO (GET /animales/{{animal_id}} debe dar 404)
GET http://localhost:8080/animales/{{animal_id}}
# aca no se le pone assert porque http 404 ya es una aserción
HTTP 404

# LIMPIEZA: ELIMINAR USUARIO DUEÑO, para que quede consistente la base 
DELETE http://localhost:8080/usuario/{{vendedor_id}}
HTTP 204



# ------ Chequeamos cada endpoint de las publicaciones ------

# CREAR VENDEDOR Y ANIMAL

# Crear Usuario Vendedor
POST http://localhost:8080/usuario/
Content-Type: application/json
{
  "nombre": "Carlos",
  "apellido": "Vendedor",
  "contacto": "carlos.vendedor@email.com"
}
HTTP 201
[Captures] #capturamos el id que nos devuelve y lo ponemos en la variable
vendedor1_id: jsonpath "$.id"

# Crear Animal
POST http://localhost:8080/animales/
Content-Type: application/json
{
  "nombre": "Milo",
  "raza": "Golden Retriever",
  "id_dueno": {{vendedor1_id}},
  "precio": "300.00"
}

HTTP 201
[Captures] # capturamos el id de animal
animal1_id: jsonpath "$.id"



# Empezamos a chequear cada endpoint para la entidad de publicacion


# CREAR PUBLICACIÓN (POST /publicacion/)
POST http://localhost:8080/publicacion/
Content-Type: application/json
{
  "id_animal": {{animal1_id}},
  "id_vendedor": {{vendedor1_id}},
  "precio": "300.00"
}

HTTP 201
[Captures]
publicacion_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.id_animal" == {{animal1_id}}
jsonpath "$.id_vendedor" == {{vendedor1_id}}


# LISTAR PUBLICACIONES ACTIVAS (GET /publicacion/)
GET http://localhost:8080/publicacion/
HTTP 200


# 3. OBTENER PUBLICACIÓN POR ID (GET /publicacion/{{publicacion_id}})
GET http://localhost:8080/publicacion/{{publicacion_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{publicacion_id}}


# ACTUALIZAR PUBLICACIÓN POR ID (PUT /publicacion/{{publicacion_id}})
PUT http://localhost:8080/publicacion/{{publicacion_id}}
Content-Type: application/json
{
  "id_animal": {{animal1_id}},
  "id_vendedor": {{vendedor1_id}},
  "precio": "350.00"
}

HTTP 200
[Asserts]
jsonpath "$.precio" == "350.00"


# ELIMINAR PUBLICACIÓN POR ID (DELETE /publicacion/{{publicacion_id}})
DELETE http://localhost:8080/publicacion/{{publicacion_id}}
HTTP 204


# VERIFICAR BORRADO (GET /publicacion/{{publicacion_id}} debe dar 404)
GET http://localhost:8080/publicacion/{{publicacion_id}}
HTTP 404