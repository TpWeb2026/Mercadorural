package main

import (
	"database/sql"
	"encoding/json"
	"net/http" //este es para poder sacar el id
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
	Contacto string `json:"contacto"` //string permite recibir strings o null en JSON
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
func aNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func (h *estructuraUsuario) getUsuarioId(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL
	partesURL := strings.Split(r.URL.Path, "/") // Ej: ["", "animales", "5"]
	idStr := partesURL[2]
}

func (h *estructuraUsuario) agregarUsuario(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) actualizarUsuarioId(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) eliminarUsuarioId(w http.ResponseWriter, r *http.Request) {

}
