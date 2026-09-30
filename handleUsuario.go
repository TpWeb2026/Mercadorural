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
type estructuraUsuario struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}

type dtoUsuario struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Contacto *string `json:"contacto"` //*string permite recibir strings o null en JSON
}

func (h *estructuraUsuario) CRUDusuarios(w http.ResponseWriter, r *http.Request) {
	cantidadURl := strings.Split(r.URL.Path, "/") // cuento la cantidad de barras que hay en la r.url.path

	if len(cantidadURl) > 3 { // chequeo que el path solo tenga involucrado 2 path, ejemplo usuario/2, si tiene 3 path lo capturamos aca
		http.Error(w, "Url no permitida", 400)
		return
	}
	// Este if lo que chequea es si tiene la longitud para que pueda tener id, si lo tiene, entra y despues se fija que metodo tiene que ejecutar
	if len(cantidadURl) == 3 && cantidadURl[2] != "" {
		switch r.Method {
		case http.MethodGet:
			h.getUsuarioId(w, r)
		case http.MethodPut:
			h.actualizarUsuarioId(w, r)
		case http.MethodDelete:
			h.eliminarUsuarioId(w, r)
		default:
			http.Error(w, "Metodo no aceptado", 405)
		}
		return //esto es para que no siga bajando
	}

	// si no entro a ninguno de los dos if, es porque solo busca lo que no tiene id
	// ejemplo ["", "usuario"] -> len = 2
	switch r.Method {
	case http.MethodGet:
		h.listaDeUsuario(w, r)
	case http.MethodPost:
		h.agregarUsuario(w, r)
	default:
		http.Error(w, "Metodo no aceptado", 405)
	}

}

// funcion que me lista todos los usuarios
func (h *estructuraUsuario) listaDeUsuario(w http.ResponseWriter, r *http.Request) {
	listaUsuario, err := h.Consultas.ListUsuarios(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//si no tiene nada, creamos una vacia para poder hacerla en formato json
	if listaUsuario == nil {
		listaUsuario = []db.Usuario{}
	}

	w.Header().Set("Content-Type", "application/json")

	//aca controlamos por si la conexion falla, para que no quede en el aire
	err = json.NewEncoder(w).Encode(listaUsuario) // listaUsuario es una lista de animales y la devolvemos en formato json por la respuesta
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Función helper para convertir string a sql.NullString
func convertidorDeString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func (h *estructuraUsuario) getUsuarioId(w http.ResponseWriter, r *http.Request) {
	// aca sacamos el id que nos llega por la url
	partesURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(partesURL[2], 10, 64)
	//chequeamos que el id sea valido
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	//buscamos ese usuario en la base de datos con la funcion de getUsuario que nos genero sqlc
	usuario, err := h.Consultas.GetUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No se encontró el usuario con ese ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Error del servidor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//devolvemos ese usuario en formato json 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuario)
}

func (h *estructuraUsuario) agregarUsuario(w http.ResponseWriter, r *http.Request) {
	//creamos el usuario de tipo dto para poder manejar el tema del string que puede ser nulo 
	var datos dtoUsuario
	err := json.NewDecoder(r.Body).Decode(&datos)
	//chequeamos que se haya decidificado correctamente
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	//reglas de negocio

	if datos.Nombre == "" {
		http.Error(w, "El nombre no puede ser vacio", http.StatusBadRequest)
		return
	}

	if datos.Apellido == "" {
		http.Error(w, "El apellido no puede ser vacio", http.StatusBadRequest)
		return
	}

	// si las reglas de negocio estan bien, que los formatos son los correctos, se puede agregar ese usuario


	//creamos una varibale de tipo db.usuario para poder agregarlo a la base de datos
	nuevo := db.CreateUsuarioParams{
		Nombre:   datos.Nombre,
		Apellido: datos.Apellido,
		Contacto: convertidorDeString(datos.Contacto),
	}

	// llamamos a la funcion que me crea el usuario dentro de la base de datos
	nuevoUsuario, err := h.Consultas.CreateUsuario(r.Context(), nuevo)

	//validamos el error
	if err != nil {
		http.Error(w, "Error al crear usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//se lo comunicamos al front de tipo json y con el cartel de que se creo correctamente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201 Created
	json.NewEncoder(w).Encode(nuevoUsuario)
}

func (h *estructuraUsuario) actualizarUsuarioId(w http.ResponseWriter, r *http.Request) {
	// primero, sacamos el id que viene desde la url
	partesURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(partesURL[2], 10, 64)
	//chequeamos que el id que nos pasaron sea correcto

	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	// creamos la variable 
	var datos dtoUsuario
	err = json.NewDecoder(r.Body).Decode(&datos)

	if  err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	//reglas de negocio

	if datos.Nombre == "" {
		http.Error(w, "El nombre no puede ser vacio", http.StatusBadRequest)
		return
	}

	if datos.Apellido == "" {
		http.Error(w, "El apellido no puede ser vacio", http.StatusBadRequest)
		return
	}

	//creamos el usuario a agregar a la base de datos
	actualizarUsu := db.UpdateUsuarioParams{
		ID:       id,
		Nombre:   datos.Nombre,
		Apellido: datos.Apellido,
		Contacto: convertidorDeString(datos.Contacto), // convertidor de string a sql.nullString para poder ponerlo en la base de datos 
	}

	usuarioActualizado, err := h.Consultas.UpdateUsuario(r.Context(), actualizarUsu) 
	//chequeamos los errores 
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No se encontró el usuario con ese ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al actualizar usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	

	// si todo esta bien 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usuarioActualizado)
}

func (h *estructuraUsuario) eliminarUsuarioId(w http.ResponseWriter, r *http.Request) {
	
	partesURL := strings.Split(r.URL.Path, "/")
	id, err := strconv.ParseInt(partesURL[2], 10, 64)
	if err != nil {
		http.Error(w, "El ID debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Verificamos que el usuario exista
	_, err = h.Consultas.GetUsuario(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "El usuario no existe", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Eliminamos
	err = h.Consultas.DeleteUsuario(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}
