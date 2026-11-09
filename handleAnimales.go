package main

import (
	"log"
	"net/http"
	"strings"
	db "tp2/db/sqlc"
	"tp2/views"
)

func (app *application) crearAnimal(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	form := views.AnimalFormData{
		Nombre: strings.TrimSpace(r.PostFormValue("nombre")),
		Raza:   strings.TrimSpace(r.PostFormValue("raza")),
		Precio: strings.TrimSpace(r.PostFormValue("precio")),
	}

	data, err := app.buildPageData(r.Context())
	if err != nil {
		log.Printf("cargar datos para crear animal: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.AnimalForm = form

	if form.Nombre == "" || form.Raza == "" {
		data.Error = "El nombre y la raza son obligatorios."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	idDueno, err := parseID(r.PostFormValue("id_dueno"))
	if err != nil {
		data.Error = "Seleccioná un dueño válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	form.IDDueno = idDueno
	data.AnimalForm = form

	precio, err := parsePrecio(form.Precio)
	if err != nil {
		data.Error = "Ingresá un precio válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	if _, err := app.consultas.GetUsuario(r.Context(), idDueno); err != nil {
		data.Error = "El dueño seleccionado no existe o está inactivo."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	_, err = app.consultas.CreateAnimal(r.Context(), db.CreateAnimalParams{
		Nombre:  form.Nombre,
		Raza:    form.Raza,
		IDDueno: idDueno,
		Precio:  precio,
	})
	if err != nil {
		log.Printf("crear animal: %v", err)
		data.Error = "No se pudo crear el animal."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#animales", http.StatusSeeOther)
}

func (app *application) eliminarAnimal(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID del animal no es válido.", http.StatusBadRequest)
		return
	}

	if err := app.consultas.DeleteAnimal(r.Context(), id); err != nil {
		log.Printf("eliminar animal %d: %v", id, err)
		if esErrorClaveForanea(err) {
			app.renderPageError(w, r, "No se puede eliminar el animal porque tiene ventas asociadas.", http.StatusConflict)
			return
		}
		app.renderPageError(w, r, "No se pudo eliminar el animal.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#animales", http.StatusSeeOther)
}
