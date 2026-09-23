package main

import (
	"encoding/json"
	"net/http"
	"strconv" //este es para poder sacar el id
	"strings"
	db "tp2/db/sqlc"
)

// esta estrucuta se hizo para poder comunicar y poder usar los metodos que nos creo sqlc
type estructuraAnimal struct {
	Consultas *db.Queries
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
	}
	//validamos que el precio sea menor a 0
	if NuevoAnimal.Precio == "" { // se hace asi porque sqlc genera los de tipo numeric como string
		http.Error(w, "El campo de precio es vacio", 400)
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

// la funcion que dado un id pasado por la url, lo busque y se fije si existe en la base de datos ese animal o no
func (h *estructuraAnimal) getAnimalId(w http.ResponseWriter, r *http.Request) {
	cantidad := strings.Split(r.URL.Path, "/")

	idUsuario, err := strconv.Atoi(cantidad[2]) // me quedo con lo que tengo en el ultimo path animales/2

	if err != nil { // vemos si no da error
		http.Error(w, "No se puede sacar el id", 400)
		return
	}


	obtenerLista, err := h.Consultas.ListAnimales(r.Context()) // r.context(), es para conectar seguro con la base de datos, si el usuario al final cencela eso, se cae la consulta de la base de datos
	if err != nil {
		http.Error(w, "NO hay una lista valida"+err.Error(), 400)
		return
	}
	//buscamos en esa lista si existe el id, range entrega dos cosas, el indice como primer parametro y el segund es el elemento
	for _, i := range obtenerLista { 
		if i.ID == int64(idUsuario) { //lo convierto a int64 para que go no me marque error de compilacion
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(i) // devolvemos en formato json el usuario que se encontro
			return
		}
	}
	//es porque no lo encotro
	http.Error(w, "No se encontro el animal"+err.Error(), 404)
}

// nos queda hacer la funcion de delete y la funcion de update con id, y ademas de eso nos queda el handle principal para llamar a cada uno d estos metodos



