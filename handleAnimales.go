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
type estructuraAnimal struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}

func (h *estructuraAnimal) CRUDanimal(w http.ResponseWriter, r *http.Request) {
	cantidadURl := strings.Split(r.URL.Path, "/") // cuento la cantidad de barras que hay en la r.url.path

	if len(cantidadURl) > 3 { // chequeo que el path solo tenga involucrado 2 path, ejemplo animales/2, si tiene 3 path lo capturamos aca
		http.Error(w, "Url no permitida", http.StatusBadRequest)
		return
	}
	// Este if lo que chequea es si tiene la longitud para que pueda tener id, si lo tiene, entra y despues se fija que metodo tiene que ejecutar
	if len(cantidadURl) == 3 && cantidadURl[2] != "" {
		switch r.Method {
		case http.MethodGet:
			h.getAnimalId(w, r)
		case http.MethodPut:
			h.actualizarAnimalID(w, r)
		case http.MethodDelete:
			h.eliminarAnimalID(w, r)
		default:
			http.Error(w, "Metodo no aceptado", http.StatusMethodNotAllowed)
		}
		return //esto es para que no siga bajando
	}

	// si no entro a ninguno de los dos if, es porque solo busca lo que no tiene id
	// ejemplo ["", "animales"] -> len = 2
	switch r.Method {
	case http.MethodGet:
		h.listaDeAnimales(w, r)
	case http.MethodPost:
		h.agregarAnimal(w, r)
	default:
		http.Error(w, "Metodo no aceptado", http.StatusMethodNotAllowed)
	}
}

// esta funcion lo que hace es listar los animales que tenemos en la base de datos, si no hay ningun animal, devuelve una lista sin nada
func (h *estructuraAnimal) listaDeAnimales(w http.ResponseWriter, r *http.Request) {
	listaAnimal, err := h.Consultas.ListAnimales(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if listaAnimal == nil {
		listaAnimal = []db.Animal{}
	}

	w.Header().Set("Content-Type", "application/json")

	//aca controlamos por si la conexion falla, para que no quede en el aire
	err = json.NewEncoder(w).Encode(listaAnimal) // listaAnimal es una lista de animales y la devolvemos en formato json por la respuesta
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// ahora creamos la otra funcion que lo que hace es hacer el post
func (h *estructuraAnimal) agregarAnimal(w http.ResponseWriter, r *http.Request) {
	var NuevoAnimal db.Animal
	//aca lo que hacemos es decodificar lo que nos mando el cliente por la solicitud
	err := json.NewDecoder(r.Body).Decode(&NuevoAnimal)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ACA TENEMOS QUE VALIDAR QUE LOS ATRIBUTOS DEL ANIMAL SEAN VALIDOS, ESTO SERIA LAS REGLAS DE NEGOCIO
	if NuevoAnimal.Nombre == "" {
		http.Error(w, "El campo de nombre es vacio", http.StatusBadRequest)
		return
	}
	//validamos que la raza no sea vacia
	if NuevoAnimal.Raza == "" {
		http.Error(w, "El campo de raza es vacio", http.StatusBadRequest)
		return
	}
	//validamos que el id dueño no sea negativo
	if NuevoAnimal.IDDueno < 0 {
		http.Error(w, "El campo de id dueño es menos a 0", http.StatusBadRequest)
		return
	}
	//validamos que el precio sea menor a 0
	if NuevoAnimal.Precio == "" { // se hace asi porque sqlc genera los de tipo numeric como string
		http.Error(w, "El campo de precio es vacio", http.StatusBadRequest)
		return
	}

	//Validamos que el dueño exista
	if _, err := h.Consultas.GetUsuario(r.Context(), NuevoAnimal.IDDueno); err != nil {
		http.Error(w, "El Dueño no existe", http.StatusBadRequest)
		return
	}

	// mapeamos los datos a los parametros que genero sqlc
	animalBase := db.CreateAnimalParams{
		Nombre:  NuevoAnimal.Nombre,
		Raza:    NuevoAnimal.Raza,
		IDDueno: NuevoAnimal.IDDueno,
		Precio:  NuevoAnimal.Precio,
	}

	// aca lo que hacemos es ejecutar la insercion de se animal nuevo a la base de datos
	animalnuevo, err := h.Consultas.CreateAnimal(r.Context(), animalBase)

	//chequemos para ver si se pudo completar la insercion o no
	if err != nil {
		http.Error(w, "NO se pudo guardar el animal en la base de datos"+err.Error(), http.StatusInternalServerError)
		return
	}

	//ahora le respondemos al cliente con el http de animal creado y el nuevo objeto
	w.Header().Set("Content-Type", "application/json")

	//le notificamos que se creo bien con el codigo 201
	w.WriteHeader(http.StatusCreated)

	// lo devolvemos en formato json
	json.NewEncoder(w).Encode(animalnuevo)
}

// devuelve un animal en formato json segun su id
func (h *estructuraAnimal) getAnimalId(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL
	partesURL := strings.Split(r.URL.Path, "/") // Ej: ["", "animales", "5"]
	idStr := partesURL[2]

	// Convertir el ID de string a int64 (ya que UpdateAnimalParams pide int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	animalObtenido, err := h.Consultas.GetAnimal(r.Context(), int64(id))

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "El animal no existe", http.StatusNotFound)
			return
		}
		http.Error(w, "Error del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// lo devolvemos en formato json
	json.NewEncoder(w).Encode(animalObtenido)
}

// esta funcion lo que hace, es actualizar un animal segun su id
func (h *estructuraAnimal) actualizarAnimalID(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL
	partesUrl := strings.Split(r.URL.Path, "/") // Ej: ["", "animales", "5"]
	idStr := partesUrl[2]                       // me quedo con la segunda parte de la url /animales/2

	// Convertir el ID de string a int64 (ya que UpdateAnimalParams pide int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	// creo una variable que tenga el struct del animal
	var datos struct {
		Nombre  string `json:"nombre"`
		Raza    string `json:"raza"`
		IDDueno int64  `json:"id_dueno"`
		Precio  string `json:"precio"`
	}

	//chequeo si hay error al pasarlo al formato json
	err = json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		http.Error(w, "JSON inválido o mal formado", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(datos.Precio)
	if err != nil {
		http.Error(w, "el precio tiene que ser un valor numerico positivo", http.StatusBadRequest)
		return
	}
	if i < 0 {
		http.Error(w, "el precio no puede ser negativo", http.StatusBadRequest)
		return
	}

	var animal db.Animal
	animal, err = h.Consultas.GetAnimal(r.Context(), id)

	if err != nil {
		http.Error(w, "el animal no se encontro", http.StatusBadRequest)
		return
	}

	if animal.IDDueno != datos.IDDueno {
		http.Error(w, "el id dueño no puede cambiar", http.StatusBadRequest)
		return
	}

	actualizacionAnimal := db.UpdateAnimalParams{
		ID:      id, // el id que saque desde la url
		Nombre:  datos.Nombre,
		Raza:    datos.Raza,
		IDDueno: datos.IDDueno,
		Precio:  datos.Precio,
	}

	//actualizacion del animal
	// el update generado por sqlc por dentras chequea que el id exista y que la clave foranea tambien exista en la tabla
	// si el dueño no existe, Postgresql aborta la transaccion inmediatamente y le devuelve a go un error de la base de datos con el estado 23503
	animalActualizado, err := h.Consultas.UpdateAnimal(r.Context(), actualizacionAnimal)

	// si el id no existe, postgres no actualiza nada y devuelve un conjunto vacio

	if err != nil {
		// Si el ID no existe en la base de datos, Scan() devuelve sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No se encontró el animal con ese ID", http.StatusNotFound)
			return
		}

		// Para cualquier otro error de base de datos
		http.Error(w, "Error al actualizar el animal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Código 200
	json.NewEncoder(w).Encode(animalActualizado)
}

func (h *estructuraAnimal) eliminarAnimalID(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID desde la URL
	partesUrl := strings.Split(r.URL.Path, "/") // Ej: ["", "animales", "5"]
	idStr := partesUrl[2]                       // me quedo con la segunda parte de la url /animales/2

	// Convertir el ID de string a int64 (ya que UpdateAnimalParams pide int64)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	//antes de borrarlo, tengo que verificar que exista ese animal
	_, err = h.Consultas.GetAnimal(r.Context(), int64(id))

	if err != nil {
		if err == sql.ErrNoRows { // si el animal no existe, el getAnimal devuelve un errnoROWN
			http.Error(w, "El animal no existe", http.StatusNotFound)
			return
		}
		http.Error(w, "Error del servidor", http.StatusInternalServerError)
		return
	}

	err = h.Consultas.DeleteAnimal(r.Context(), id)
	if err != nil {
		http.Error(w, "No se pudo borrar el animal", http.StatusInternalServerError)
		return
	}
	// se elimino correctamente
	w.WriteHeader(http.StatusNoContent)
}
