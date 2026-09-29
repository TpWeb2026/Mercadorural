package main

import (
	"database/sql"
	"encoding/json"
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
		http.Error(w, "Url no permitida", 400)
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
			http.Error(w, "Metodo no aceptado", 405)
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
		http.Error(w, "Metodo no aceptado", 405)
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
		http.Error(w, err.Error(), 400)
		return
	}

	// ACA TENEMOS QUE VALIDAR QUE LOS ATRIBUTOS DEL ANIMAL SEAN VALIDOS, ESTO SERIA LAS REGLAS DE NEGOCIO
	if NuevoAnimal.Nombre == "" {
		http.Error(w, "El campo de nombre es vacio", 400)
		return
	}
	//validamos que la raza no sea vacia
	if NuevoAnimal.Raza == "" {
		http.Error(w, "El campo de raza es vacio", 400)
		return
	}
	//validamos que el id dueño no sea negativo
	if NuevoAnimal.IDDueno < 0 {
		http.Error(w, "El campo de id dueño es menos a 0", 400)
		return
	}
	//validamos que el precio sea menor a 0
	if NuevoAnimal.Precio == "" { // se hace asi porque sqlc genera los de tipo numeric como string
		http.Error(w, "El campo de precio es vacio", 400)
		return
	}

	//Validamos que el dueño exista
	if _, err := h.Consultas.GetUsuario(r.Context(), NuevoAnimal.IDDueno); err != nil {
		http.Error(w, "El Dueño no existe", 400)
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
	w.WriteHeader(201)

	// lo devolvemos en formato json
	json.NewEncoder(w).Encode(animalnuevo)
}

func (h *estructuraAnimal) getAnimalId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Falta el parámetro 'id'", 400)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "El id debe ser un número válido", 400)
		return
	}

	animalObtenido, err := h.Consultas.GetAnimal(r.Context(), int64(id))

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "El animal no existe", 404)
			return
		}
		http.Error(w, "Error del servidor", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	/*
		if err := json.NewEncoder(w).Encode(animalObtenido); err != nil {
			// Solo logueamos el error interno, ya que el header HTTP ya fue enviado al cliente
			// log.Printf("Error serializando JSON de animal: %v", err)
		}
	*/
}

func (h *estructuraAnimal) actualizarAnimalID(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraAnimal) eliminarAnimalID(w http.ResponseWriter, r *http.Request) {

}

// nos queda hacer la funcion de delete y la funcion de update con id, y ademas de eso nos queda el handle principal para llamar a cada uno d estos metodos
