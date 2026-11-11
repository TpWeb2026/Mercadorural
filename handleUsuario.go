package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	db "tp2/db/sqlc"
	"tp2/views"
)

func (app *application) seleccionarOperador(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	id, err := parseID(r.PostFormValue("id_usuario"))
	if err != nil {
		app.renderPageError(w, r, "Seleccioná un usuario válido.", http.StatusBadRequest)
		return
	}

	if _, err := app.consultas.GetUsuario(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "El usuario seleccionado no existe o está inactivo.", http.StatusBadRequest)
			return
		}
		log.Printf("buscar operador %d: %v", id, err)
		app.renderPageError(w, r, "No se pudo verificar el usuario.", http.StatusInternalServerError)
		return
	}

	guardarOperadorCookie(w, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) nuevoUsuarioForm(w http.ResponseWriter, r *http.Request) {
	app.renderNuevoUsuario(w, r, views.UsuarioFormData{}, "", http.StatusOK)
}

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

	if form.Nombre == "" || form.Apellido == "" {
		app.renderNuevoUsuario(w, r, form, "El nombre y el apellido son obligatorios.", http.StatusBadRequest)
		return
	}

	usuario, err := app.consultas.CreateUsuario(r.Context(), db.CreateUsuarioParams{
		Nombre:   form.Nombre,
		Apellido: form.Apellido,
		Contacto: convertirContacto(form.Contacto),
	})
	if err != nil {
		log.Printf("crear usuario: %v", err)
		app.renderNuevoUsuario(w, r, form, "No se pudo crear el usuario.", http.StatusInternalServerError)
		return
	}

	guardarOperadorCookie(w, usuario.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) eliminarUsuario(w http.ResponseWriter, r *http.Request) {
	operadorID, ok := operadorDeCookie(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de eliminar la cuenta.", http.StatusBadRequest)
		return
	}

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID del usuario no es válido.", http.StatusBadRequest)
		return
	}

	if id != operadorID {
		app.renderPageError(w, r, "Solo podés eliminar tu propia cuenta.", http.StatusForbidden)
		return
	}

	if _, err := app.consultas.GetUsuario(r.Context(), id); err != nil {
		app.renderPageError(w, r, "El usuario no existe o ya está inactivo.", http.StatusNotFound)
		return
	}

	tx, err := app.conexion.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("iniciar transacción de baja de usuario: %v", err)
		app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	consultas := app.consultas.WithTx(tx)

	publicaciones, err := consultas.ListPublicacionesActivasPorVendedor(r.Context(), id)
	if err != nil {
		log.Printf("listar publicaciones activas del usuario %d: %v", id, err)
		app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
		return
	}

	for _, publicacion := range publicaciones {
		if err := consultas.DeletePublicacion(r.Context(), publicacion.ID); err != nil {
			log.Printf("eliminar publicación %d del usuario %d: %v", publicacion.ID, id, err)
			app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
			return
		}

		if _, err := consultas.DeleteAnimalSiSinPublicaciones(r.Context(), publicacion.IDAnimal); err != nil {
			log.Printf("eliminar animal %d de la publicación %d: %v", publicacion.IDAnimal, publicacion.ID, err)
			app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
			return
		}
	}

	if err := consultas.DeleteUsuario(r.Context(), id); err != nil {
		log.Printf("eliminar usuario %d: %v", id, err)
		app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("confirmar baja del usuario %d: %v", id, err)
		app.renderPageError(w, r, "No se pudo eliminar la cuenta.", http.StatusInternalServerError)
		return
	}

	limpiarOperadorCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
