package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	db "tp2/db/sqlc"
	"tp2/views"

	"github.com/a-h/templ"
)

const operadorCookie = "operador_id"

type application struct {
	consultas *db.Queries
	conexion  *sql.DB
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.index)
	mux.HandleFunc("POST /operador", app.seleccionarOperador)
	mux.HandleFunc("GET /usuario/nuevo", app.nuevoUsuarioForm)
	mux.HandleFunc("POST /usuario", app.crearUsuario)
	mux.HandleFunc("POST /usuario/{id}/eliminar", app.eliminarUsuario)
	mux.HandleFunc("POST /publicacion", app.crearPublicacion)
	mux.HandleFunc("GET /publicacion/{id}/editar", app.editarPublicacionForm)
	mux.HandleFunc("POST /publicacion/{id}/editar", app.actualizarPublicacion)
	mux.HandleFunc("POST /publicacion/{id}/comprar", app.comprarPublicacion)
	mux.HandleFunc("POST /publicacion/{id}/eliminar", app.eliminarPublicacion)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return mux
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	operadorID, _ := operadorDeCookie(r)

	data, err := app.buildPageData(r.Context(), operadorID)
	if err != nil {
		log.Printf("cargar página principal: %v", err)
		http.Error(w, "No se pudieron cargar los datos", http.StatusInternalServerError)
		return
	}

	if operadorID > 0 && data.Operador == nil {
		limpiarOperadorCookie(w)
	}

	app.renderIndex(w, r, data, http.StatusOK)
}

func (app *application) buildPageData(ctx context.Context, operadorID int64) (views.PageData, error) {
	usuarios, err := app.consultas.ListUsuarios(ctx)
	if err != nil {
		return views.PageData{}, fmt.Errorf("listar usuarios: %w", err)
	}

	animales, err := app.consultas.ListAnimales(ctx)
	if err != nil {
		return views.PageData{}, fmt.Errorf("listar animales: %w", err)
	}

	publicaciones, err := app.consultas.ListPublicaciones(ctx)
	if err != nil {
		return views.PageData{}, fmt.Errorf("listar publicaciones: %w", err)
	}

	data := views.PageData{}
	usuariosPorID := make(map[int64]db.Usuario, len(usuarios))

	for _, usuario := range usuarios {
		usuariosPorID[usuario.ID] = usuario
		nombre := views.NombreCompleto(usuario.Nombre, usuario.Apellido)
		data.Usuarios = append(data.Usuarios, views.Opcion{
			ID:    usuario.ID,
			Texto: nombre,
		})
	}

	if operador, ok := usuariosPorID[operadorID]; ok {
		data.Operador = &views.OperadorVista{
			ID:             operador.ID,
			NombreCompleto: views.NombreCompleto(operador.Nombre, operador.Apellido),
			Contacto:       views.FormatearContacto(operador.Contacto),
		}
	}

	animalesPorID := make(map[int64]db.Animal, len(animales))

	for _, animal := range animales {
		animalesPorID[animal.ID] = animal
	}

	for _, publicacion := range publicaciones {
		animalNombre := fmt.Sprintf("Animal #%d", publicacion.IDAnimal)
		raza := "Sin datos"
		if animal, ok := animalesPorID[publicacion.IDAnimal]; ok {
			animalNombre = animal.Nombre
			raza = animal.Raza
		}

		vendedor := fmt.Sprintf("Usuario #%d", publicacion.IDVendedor)
		if usuario, ok := usuariosPorID[publicacion.IDVendedor]; ok {
			vendedor = views.NombreCompleto(usuario.Nombre, usuario.Apellido)
		}

		vista := views.PublicacionVista{
			ID:       publicacion.ID,
			Animal:   animalNombre,
			Raza:     raza,
			Vendedor: vendedor,
			Precio:   views.FormatearPrecio(publicacion.Precio),
			Activa:   publicacion.Activo,
		}

		if operadorID > 0 && publicacion.IDVendedor == operadorID {
			data.PublicacionesPropias = append(data.PublicacionesPropias, vista)
		} else {
			data.PublicacionesAjenas = append(data.PublicacionesAjenas, vista)
		}
	}

	return data, nil
}

func (app *application) renderComponent(w http.ResponseWriter, r *http.Request, componente templ.Component, status int) {
	var buffer bytes.Buffer

	if err := componente.Render(r.Context(), &buffer); err != nil {
		log.Printf("renderizar componente: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("escribir respuesta: %v", err)
	}
}

func (app *application) renderIndex(w http.ResponseWriter, r *http.Request, data views.PageData, status int) {
	app.renderComponent(w, r, views.IndexPage(data), status)
}

func (app *application) renderPageError(w http.ResponseWriter, r *http.Request, mensaje string, status int) {
	operadorID, _ := operadorDeCookie(r)

	data, err := app.buildPageData(r.Context(), operadorID)
	if err != nil {
		log.Printf("cargar datos para error: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.Error = mensaje
	app.renderIndex(w, r, data, status)
}

func (app *application) renderNuevoUsuario(w http.ResponseWriter, r *http.Request, form views.UsuarioFormData, mensaje string, status int) {
	app.renderComponent(w, r, views.NuevoUsuarioPage(form, mensaje), status)
}

func (app *application) renderEditarPublicacion(w http.ResponseWriter, r *http.Request, operador *views.OperadorVista, form views.EdicionPublicacionData, mensaje string, status int) {
	app.renderComponent(w, r, views.EditarPublicacionPage(operador, form, mensaje), status)
}

func (app *application) operadorVista(ctx context.Context, id int64) (*views.OperadorVista, error) {
	operador, err := app.consultas.GetUsuario(ctx, id)
	if err != nil {
		return nil, err
	}

	return &views.OperadorVista{
		ID:             operador.ID,
		NombreCompleto: views.NombreCompleto(operador.Nombre, operador.Apellido),
		Contacto:       views.FormatearContacto(operador.Contacto),
	}, nil
}

func operadorDeCookie(r *http.Request) (int64, bool) {
	cookie, err := r.Cookie(operadorCookie)
	if err != nil {
		return 0, false
	}

	id, err := parseID(cookie.Value)
	if err != nil {
		return 0, false
	}

	return id, true
}

func (app *application) operadorActivo(r *http.Request) (int64, bool) {
	id, ok := operadorDeCookie(r)
	if !ok {
		return 0, false
	}

	if _, err := app.consultas.GetUsuario(r.Context(), id); err != nil {
		return 0, false
	}

	return id, true
}

func guardarOperadorCookie(w http.ResponseWriter, id int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     operadorCookie,
		Value:    strconv.FormatInt(id, 10),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func limpiarOperadorCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     operadorCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
