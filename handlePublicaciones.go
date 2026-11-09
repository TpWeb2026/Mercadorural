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

func (app *application) crearPublicacion(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	form := views.PublicacionFormData{
		Precio: strings.TrimSpace(r.PostFormValue("precio")),
	}

	data, err := app.buildPageData(r.Context())
	if err != nil {
		log.Printf("cargar datos para crear publicación: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.PublicacionForm = form

	idAnimal, err := parseID(r.PostFormValue("id_animal"))
	if err != nil {
		data.Error = "Seleccioná un animal válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	form.IDAnimal = idAnimal
	data.PublicacionForm = form

	idVendedor, err := parseID(r.PostFormValue("id_vendedor"))
	if err != nil {
		data.Error = "Seleccioná un vendedor válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	form.IDVendedor = idVendedor
	data.PublicacionForm = form

	precio, err := parsePrecio(form.Precio)
	if err != nil {
		data.Error = "Ingresá un precio válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	if _, err := app.consultas.GetAnimal(r.Context(), idAnimal); err != nil {
		data.Error = "El animal seleccionado no existe."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	if _, err := app.consultas.GetUsuario(r.Context(), idVendedor); err != nil {
		data.Error = "El vendedor seleccionado no existe o está inactivo."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	_, err = app.consultas.CreatePublicacion(r.Context(), db.CreatePublicacionParams{
		IDAnimal:   idAnimal,
		Precio:     precio,
		IDVendedor: idVendedor,
	})
	if err != nil {
		log.Printf("crear publicación: %v", err)
		data.Error = "No se pudo crear la publicación."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#publicaciones", http.StatusSeeOther)
}

func (app *application) comprarPublicacion(w http.ResponseWriter, r *http.Request) {
	idPublicacion, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	idComprador, err := parseID(r.PostFormValue("id_comprador"))
	if err != nil {
		app.renderPageError(w, r, "Seleccioná un comprador válido.", http.StatusBadRequest)
		return
	}

	if _, err := app.consultas.GetUsuario(r.Context(), idComprador); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "El comprador no existe o está inactivo.", http.StatusBadRequest)
			return
		}
		log.Printf("buscar comprador %d: %v", idComprador, err)
		app.renderPageError(w, r, "No se pudo verificar el comprador.", http.StatusInternalServerError)
		return
	}

	tx, err := app.conexion.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("iniciar transacción de compra: %v", err)
		app.renderPageError(w, r, "No se pudo iniciar la compra.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	consultas := app.consultas.WithTx(tx)

	publicacion, err := consultas.GetPublicacionParaCompra(r.Context(), idPublicacion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "La publicación no existe.", http.StatusNotFound)
			return
		}
		log.Printf("buscar publicación %d para compra: %v", idPublicacion, err)
		app.renderPageError(w, r, "No se pudo verificar la publicación.", http.StatusInternalServerError)
		return
	}

	if !publicacion.Activo {
		app.renderPageError(w, r, "La publicación ya fue vendida.", http.StatusConflict)
		return
	}

	if publicacion.IDVendedor == idComprador {
		app.renderPageError(w, r, "El comprador no puede ser el mismo que el vendedor.", http.StatusBadRequest)
		return
	}

	filas, err := consultas.MarcarPublicacionVendida(r.Context(), idPublicacion)
	if err != nil {
		log.Printf("marcar publicación %d como vendida: %v", idPublicacion, err)
		app.renderPageError(w, r, "No se pudo registrar la venta.", http.StatusInternalServerError)
		return
	}

	if filas == 0 {
		app.renderPageError(w, r, "La publicación ya fue vendida.", http.StatusConflict)
		return
	}

	_, err = consultas.CreateVenta(r.Context(), db.CreateVentaParams{
		IDPublicacion: idPublicacion,
		IDVendedor:    publicacion.IDVendedor,
		IDComprador:   idComprador,
	})
	if err != nil {
		log.Printf("crear venta de publicación %d: %v", idPublicacion, err)
		app.renderPageError(w, r, "No se pudo registrar la venta.", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("confirmar venta de publicación %d: %v", idPublicacion, err)
		app.renderPageError(w, r, "No se pudo confirmar la venta.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#publicaciones", http.StatusSeeOther)
}

func (app *application) eliminarPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
		return
	}

	if err := app.consultas.DeletePublicacion(r.Context(), id); err != nil {
		log.Printf("eliminar publicación %d: %v", id, err)
		if esErrorClaveForanea(err) {
			app.renderPageError(w, r, "No se puede eliminar la publicación porque tiene una venta asociada.", http.StatusConflict)
			return
		}
		app.renderPageError(w, r, "No se pudo eliminar la publicación.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#publicaciones", http.StatusSeeOther)
}
