package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http" //este es para poder sacar el id
	"strconv"
	"strings"
	db "tp2/db/sqlc"
)

// esta estrucuta se hizo para poder comunicar y poder usar los metodos que nos creo sqlc
type estructuraPublicacion struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}

func (h *estructuraPublicacion) CRUDpublicacion(w http.ResponseWriter, r *http.Request) {
	cantidadURL := strings.Split(r.URL.Path, "/")

	if len(cantidadURL) > 3 {
		http.Error(w, "URL no permitida", http.StatusBadRequest)
		return
	}

	// CASO CON ID: /publicaciones/5
	if len(cantidadURL) == 3 && cantidadURL[2] != "" {
		switch r.Method {
		case http.MethodGet:
			h.getPublicacionId(w, r)
		case http.MethodPut:
			h.actualizarPublicacionId(w, r)
		case http.MethodDelete:
			h.eliminarPublicacionId(w, r)
		default:
			http.Error(w, "Método no aceptado", http.StatusMethodNotAllowed)
		}
		return
	}

	// CASO SIN ID: /publicaciones
	switch r.Method {
	case http.MethodGet:
		h.listaDePublicaciones(w, r)
	case http.MethodPost:
		h.agregarPublicacion(w, r)
	default:
		http.Error(w, "Método no aceptado", http.StatusMethodNotAllowed)
	}
}

func (h *estructuraPublicacion) listaDePublicaciones(w http.ResponseWriter, r *http.Request) {
	//guardamos en una variables todas las publicacione activas que nos devuelve la base de datos
	publicaciones, err := h.Consultas.ListPublicacionesActivas(r.Context())

	//chequeamos errores
	if err != nil {
		http.Error(w, "Error al obtener publicaciones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// si no hay ninguna publicacion, entonces creamos una vacia para poder devolvesela al usuario en formato json
	if publicaciones == nil {
		publicaciones = []db.Publicacion{}
	}

	// se la devolvemos al usuario en formato json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicaciones)

}

func (h *estructuraPublicacion) getPublicacionId(w http.ResponseWriter, r *http.Request) {
	//sacamos el url que nos viene desde el request
	partesURL := strings.Split(r.URL.Path, "/")

	//lo convertimos a int64 porque esta en string
	id, err := strconv.ParseInt(partesURL[2], 10, 64)

	//chequeamos que el id sea correcto
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	//llamamos a la funcion que nos devuelve la publicacion dependiendo de un id
	publicacion, err := h.Consultas.GetPublicacion(r.Context(), id)

	//chequeamos los errores
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No se encontró la publicación con ese ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//encapsulamos la respuesta en formato json y se la mostramos al usuario
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound) // aca se cambia el StatusOk que da el 200 por el 404 como dice el enunciado
	json.NewEncoder(w).Encode(publicacion)
}

func (h *estructuraPublicacion) actualizarPublicacionId(w http.ResponseWriter, r *http.Request) {
	// aca sacamos el id que nos llega por la url
	partesURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(partesURL[2], 10, 64)
	//chequeamos que el id sea valido
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	var datos db.Publicacion
	err = json.NewDecoder(r.Body).Decode(&datos)

	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	//convierto el precio que estaba en string a int para pdoer compararlo
	datosPrecio, err := strconv.ParseInt(datos.Precio, 10, 64)

	// Validaciones básicas de campos obligatorios, regla de negocio
	if datos.IDAnimal <= 0 || datos.IDVendedor <= 0 || datosPrecio < 0 {
		http.Error(w, "Los campos id_animal, id_vendedor y precio son obligatorios", http.StatusBadRequest)
		return
	}

	// si todo paso correctamente, podemos actualizar la publicacion

	params := db.UpdatePublicacionParams{
		IDAnimal:   datos.IDAnimal,
		Precio:     datos.Precio,
		IDVendedor: datos.IDVendedor,
	}

	//ahora traemos buscar el id vendedor y el id animal antes de actualizar la publicacion

	//traer publicacion
	// validar

	//traer usuario
	// validar

	//traer animal
	// validar

	_, err = h.Consultas.GetPublicacion(r.Context(), id)
	if err != nil {
		http.Error(w, "No existe la publicacion", http.StatusBadRequest)
		return
	}

	_, err = h.Consultas.GetUsuario(r.Context(), datos.IDVendedor)
	if err != nil {
		http.Error(w, "No existe el vendedor", http.StatusBadRequest)
		return
	}

	_, err = h.Consultas.GetAnimal(r.Context(), datos.IDAnimal)
	if err != nil {
		http.Error(w, "No existe el animal", http.StatusBadRequest)
		return
	}

	// todo anda correctamente, por lo tanto, ahora se puede actualizar la publicacion

	nuevaPublicacion, err := h.Consultas.UpdatePublicacion(r.Context(), params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(nuevaPublicacion)
}

func (h *estructuraPublicacion) agregarPublicacion(w http.ResponseWriter, r *http.Request) {
	var nuevaPublicacion db.Publicacion
	//aca lo que hacemos es decodificar lo que nos mando el cliente por la solicitud
	err := json.NewDecoder(r.Body).Decode(&nuevaPublicacion)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validamos que el precio no sea vacio para que el parse no falle al convertir
	if nuevaPublicacion.Precio == "" {
		http.Error(w, "el precio es un campo obligatorio", http.StatusBadRequest)
		return
	}

	// aca usamos el parseint para pasar el precio a numero, pero chequemos err por si trato de pasar un string ej: "abc" a numero
	precioconv, err := strconv.ParseInt(nuevaPublicacion.Precio, 10, 64)
	if err != nil {
		http.Error(w, "el formato del precio invalido", http.StatusBadRequest)
		return
	}
	// aca validamos reglas de negocio
	if nuevaPublicacion.IDVendedor < 0 || nuevaPublicacion.IDAnimal < 0 || precioconv < 0 {
		http.Error(w, "el id vendedor y/o el id animal no pueden ser menor a cero, y el precio no puede ser negativo", http.StatusBadRequest)
		return
	}
	// aca validamos que exista el animal, para no tener inconsistencia en la DB
	_, err = h.Consultas.GetAnimal(r.Context(), nuevaPublicacion.IDAnimal)
	if err != nil {
		http.Error(w, "el animal no existe en la DB", http.StatusBadRequest)
		return
	}
	// aca validamos que exista el usuario, para no tener inconsistencia en la DB
	_, err = h.Consultas.GetUsuario(r.Context(), nuevaPublicacion.IDVendedor)
	if err != nil {
		http.Error(w, "el usuario no existe en la DB", http.StatusBadRequest)
		return
	}
	//armamos los parametros para la creacion de la publicacion
	params := db.CreatePublicacionParams{
		IDAnimal:   nuevaPublicacion.IDAnimal,
		Precio:     nuevaPublicacion.Precio,
		IDVendedor: nuevaPublicacion.IDVendedor,
	}
	// guardamos la publicacion en la DB
	publicacion, err := h.Consultas.CreatePublicacion(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear la publicacion en la DB", http.StatusInternalServerError)
		return
	}

	// si todo esta bien
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // se cambio el codigo de 200 a 201 porque en el enunciado dice devolver codigo 201
	json.NewEncoder(w).Encode(publicacion)
}
func (h *estructuraPublicacion) eliminarPublicacionId(w http.ResponseWriter, r *http.Request) {
	// 	ACA SUPONGO QUE CON ESTO OBTENGO EL ID
	partesURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(partesURL[2], 10, 64)
	if err != nil {
		http.Error(w, "El ID debe ser un número válido", http.StatusBadRequest)
		return
	}

	// chequeamos si la publicacion existe o no
	_, err = h.Consultas.GetPublicacion(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "La publicación no existe", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en el servidor al buscar la publicación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// aca eliminamos la publicacion usando el metodo delete de publicacion
	err = h.Consultas.DeletePublicacion(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar la publicacion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // devuelvo estado 204
}
