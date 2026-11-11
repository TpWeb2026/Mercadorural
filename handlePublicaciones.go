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
	operadorID, ok := app.operadorActivo(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de publicar.", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	form := views.PublicacionFormData{
		NombreAnimal: strings.TrimSpace(r.PostFormValue("nombre_animal")),
		Raza:         strings.TrimSpace(r.PostFormValue("raza")),
		Precio:       strings.TrimSpace(r.PostFormValue("precio")),
	}

	data, err := app.buildPageData(r.Context(), operadorID)
	if err != nil {
		log.Printf("cargar datos para crear publicación: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.PublicacionForm = form

	if form.NombreAnimal == "" || form.Raza == "" {
		data.Error = "El nombre del animal y la raza son obligatorios."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	precio, err := parsePrecio(form.Precio)
	if err != nil {
		data.Error = "Ingresá un precio válido."
		app.renderIndex(w, r, data, http.StatusBadRequest)
		return
	}

	tx, err := app.conexion.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("iniciar transacción de publicación: %v", err)
		data.Error = "No se pudo crear la publicación."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	consultas := app.consultas.WithTx(tx)

	animal, err := consultas.CreateAnimal(r.Context(), db.CreateAnimalParams{
		Nombre:  form.NombreAnimal,
		Raza:    form.Raza,
		IDDueno: operadorID,
		Precio:  precio,
	})
	if err != nil {
		log.Printf("crear animal de publicación: %v", err)
		data.Error = "No se pudo crear la publicación."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	_, err = consultas.CreatePublicacion(r.Context(), db.CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     precio,
		IDVendedor: operadorID,
	})
	if err != nil {
		log.Printf("crear publicación del animal %d: %v", animal.ID, err)
		data.Error = "No se pudo crear la publicación."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("confirmar publicación: %v", err)
		data.Error = "No se pudo crear la publicación."
		app.renderIndex(w, r, data, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#mis-publicaciones", http.StatusSeeOther)
}

func (app *application) editarPublicacionForm(w http.ResponseWriter, r *http.Request) {
	operadorID, ok := app.operadorActivo(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de editar.", http.StatusBadRequest)
		return
	}

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
		return
	}

	publicacion, err := app.consultas.GetPublicacionParaCompra(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "La publicación no existe.", http.StatusNotFound)
			return
		}
		log.Printf("buscar publicación %d para editar: %v", id, err)
		app.renderPageError(w, r, "No se pudo cargar la publicación.", http.StatusInternalServerError)
		return
	}

	if publicacion.IDVendedor != operadorID {
		app.renderPageError(w, r, "Solo podés editar tus propias publicaciones.", http.StatusForbidden)
		return
	}

	if !publicacion.Activo {
		app.renderPageError(w, r, "No se puede editar una publicación finalizada.", http.StatusConflict)
		return
	}

	animal, err := app.consultas.GetAnimal(r.Context(), publicacion.IDAnimal)
	if err != nil {
		log.Printf("buscar animal %d de publicación %d: %v", publicacion.IDAnimal, publicacion.ID, err)
		app.renderPageError(w, r, "No se pudo cargar el animal.", http.StatusInternalServerError)
		return
	}

	operador, err := app.operadorVista(r.Context(), operadorID)
	if err != nil {
		log.Printf("cargar operador %d: %v", operadorID, err)
		app.renderPageError(w, r, "No se pudo cargar el operador.", http.StatusInternalServerError)
		return
	}

	form := views.EdicionPublicacionData{
		PublicacionID: publicacion.ID,
		IDAnimal:      publicacion.IDAnimal,
		NombreAnimal:  animal.Nombre,
		Raza:          animal.Raza,
		Precio:        publicacion.Precio,
	}

	app.renderEditarPublicacion(w, r, operador, form, "", http.StatusOK)
}

func (app *application) actualizarPublicacion(w http.ResponseWriter, r *http.Request) {
	operadorID, ok := app.operadorActivo(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de editar.", http.StatusBadRequest)
		return
	}

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
		return
	}

	publicacion, err := app.consultas.GetPublicacionParaCompra(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "La publicación no existe.", http.StatusNotFound)
			return
		}
		log.Printf("buscar publicación %d para actualizar: %v", id, err)
		app.renderPageError(w, r, "No se pudo cargar la publicación.", http.StatusInternalServerError)
		return
	}

	if publicacion.IDVendedor != operadorID {
		app.renderPageError(w, r, "Solo podés editar tus propias publicaciones.", http.StatusForbidden)
		return
	}

	if !publicacion.Activo {
		app.renderPageError(w, r, "No se puede editar una publicación finalizada.", http.StatusConflict)
		return
	}

	animal, err := app.consultas.GetAnimal(r.Context(), publicacion.IDAnimal)
	if err != nil {
		log.Printf("buscar animal %d de publicación %d: %v", publicacion.IDAnimal, publicacion.ID, err)
		app.renderPageError(w, r, "No se pudo cargar el animal.", http.StatusInternalServerError)
		return
	}

	operador, err := app.operadorVista(r.Context(), operadorID)
	if err != nil {
		log.Printf("cargar operador %d: %v", operadorID, err)
		app.renderPageError(w, r, "No se pudo cargar el operador.", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulario inválido", http.StatusBadRequest)
		return
	}

	form := views.EdicionPublicacionData{
		PublicacionID: publicacion.ID,
		IDAnimal:      publicacion.IDAnimal,
		NombreAnimal:  strings.TrimSpace(r.PostFormValue("nombre_animal")),
		Raza:          strings.TrimSpace(r.PostFormValue("raza")),
		Precio:        strings.TrimSpace(r.PostFormValue("precio")),
	}

	if form.NombreAnimal == "" || form.Raza == "" {
		app.renderEditarPublicacion(w, r, operador, form, "El nombre del animal y la raza son obligatorios.", http.StatusBadRequest)
		return
	}

	precio, err := parsePrecio(form.Precio)
	if err != nil {
		app.renderEditarPublicacion(w, r, operador, form, "Ingresá un precio válido.", http.StatusBadRequest)
		return
	}

	tx, err := app.conexion.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("iniciar transacción de edición: %v", err)
		app.renderEditarPublicacion(w, r, operador, form, "No se pudo actualizar la publicación.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	consultas := app.consultas.WithTx(tx)

	_, err = consultas.UpdateAnimal(r.Context(), db.UpdateAnimalParams{
		ID:      animal.ID,
		Nombre:  form.NombreAnimal,
		Raza:    form.Raza,
		IDDueno: animal.IDDueno,
		Precio:  precio,
	})
	if err != nil {
		log.Printf("actualizar animal %d: %v", animal.ID, err)
		app.renderEditarPublicacion(w, r, operador, form, "No se pudo actualizar la publicación.", http.StatusInternalServerError)
		return
	}

	_, err = consultas.UpdatePublicacion(r.Context(), db.UpdatePublicacionParams{
		ID:         publicacion.ID,
		IDAnimal:   publicacion.IDAnimal,
		Precio:     precio,
		IDVendedor: publicacion.IDVendedor,
	})
	if err != nil {
		log.Printf("actualizar publicación %d: %v", publicacion.ID, err)
		app.renderEditarPublicacion(w, r, operador, form, "No se pudo actualizar la publicación.", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("confirmar edición de publicación %d: %v", publicacion.ID, err)
		app.renderEditarPublicacion(w, r, operador, form, "No se pudo actualizar la publicación.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#mis-publicaciones", http.StatusSeeOther)
}

func (app *application) comprarPublicacion(w http.ResponseWriter, r *http.Request) {
	operadorID, ok := app.operadorActivo(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de comprar.", http.StatusBadRequest)
		return
	}

	idPublicacion, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
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

	if publicacion.IDVendedor == operadorID {
		app.renderPageError(w, r, "No podés comprar tu propia publicación.", http.StatusBadRequest)
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
		IDComprador:   operadorID,
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

	http.Redirect(w, r, "/#publicaciones-ajenas", http.StatusSeeOther)
}

func (app *application) eliminarPublicacion(w http.ResponseWriter, r *http.Request) {
	operadorID, ok := app.operadorActivo(r)
	if !ok {
		app.renderPageError(w, r, "Seleccioná un operador antes de eliminar.", http.StatusBadRequest)
		return
	}

	id, err := parseID(r.PathValue("id"))
	if err != nil {
		app.renderPageError(w, r, "El ID de la publicación no es válido.", http.StatusBadRequest)
		return
	}

	publicacion, err := app.consultas.GetPublicacionParaCompra(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.renderPageError(w, r, "La publicación no existe.", http.StatusNotFound)
			return
		}
		log.Printf("buscar publicación %d para eliminar: %v", id, err)
		app.renderPageError(w, r, "No se pudo cargar la publicación.", http.StatusInternalServerError)
		return
	}

	if publicacion.IDVendedor != operadorID {
		app.renderPageError(w, r, "Solo podés eliminar tus propias publicaciones.", http.StatusForbidden)
		return
	}

	if !publicacion.Activo {
		app.renderPageError(w, r, "No se puede eliminar una publicación finalizada.", http.StatusConflict)
		return
	}

	tx, err := app.conexion.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("iniciar transacción de eliminación: %v", err)
		app.renderPageError(w, r, "No se pudo eliminar la publicación.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	consultas := app.consultas.WithTx(tx)

	if err := consultas.DeletePublicacion(r.Context(), publicacion.ID); err != nil {
		log.Printf("eliminar publicación %d: %v", publicacion.ID, err)
		app.renderPageError(w, r, "No se pudo eliminar la publicación.", http.StatusInternalServerError)
		return
	}

	if _, err := consultas.DeleteAnimalSiSinPublicaciones(r.Context(), publicacion.IDAnimal); err != nil {
		log.Printf("eliminar animal %d de publicación %d: %v", publicacion.IDAnimal, publicacion.ID, err)
		app.renderPageError(w, r, "No se pudo eliminar la publicación.", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("confirmar eliminación de publicación %d: %v", publicacion.ID, err)
		app.renderPageError(w, r, "No se pudo eliminar la publicación.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/#mis-publicaciones", http.StatusSeeOther)
}
