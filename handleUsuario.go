package main

import (
	"log"
	"net/http"
	"strings"
	db "tp2/db/sqlc"
	"tp2/views"
)

func (app *application) crearUsuario(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	form := views.UsuarioFormData{
		Nombre:   strings.TrimSpace(r.PostFormValue("nombre")),
		Apellido: strings.TrimSpace(r.PostFormValue("apellido")),
		Contacto: strings.TrimSpace(r.PostFormValue("contacto")),
	}

	data, err := app.buildPageData(r.Context())
	if err != nil {
		log.Printf("cargar datos para crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.UsuarioForm = form

	if form.Nombre == "" || form.Apellido == "" {
		data.Error = "El nombre y el apellido son obligatorios."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	_, err = app.consultas.CreateUsuario(r.Context(), db.CreateUsuarioParams{
		Nombre:   form.Nombre,
		Apellido: form.Apellido,
		Contacto: convertirContacto(form.Contacto),
	})
	if err != nil {
		log.Printf("crear usuario: %v", err)
		data.Error = "No se pudo crear el usuario."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#usuarios", http.StatusSeeOther)
}

func (app *application) eliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID del usuario no es válido.", http.StatusBadRequest)
		return
	}

	if err := app.consultas.DeleteUsuario(r.Context(), id); err != nil {
		log.Printf("eliminar usuario %d: %v", id, err)
		app.renderPageError(w, r, "No se pudo eliminar el usuario.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#usuarios", http.StatusSeeOther)
}
