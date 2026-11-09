package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	db "tp2/db/sqlc"
	"tp2/views"
)

type application struct {
	consultas *db.Queries
	conexion  *sql.DB
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.index)
	mux.HandleFunc("POST /usuario", app.crearUsuario)
	mux.HandleFunc("POST /animales", app.crearAnimal)
	mux.HandleFunc("POST /publicacion", app.crearPublicacion)
	mux.HandleFunc("POST /publicacion/{id}/comprar", app.comprarPublicacion)
	mux.HandleFunc("POST /usuario/{id}/eliminar", app.eliminarUsuario)
	mux.HandleFunc("POST /animales/{id}/eliminar", app.eliminarAnimal)
	mux.HandleFunc("POST /publicacion/{id}/eliminar", app.eliminarPublicacion)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return mux
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	data, err := app.buildPageData(r.Context())
	if err != nil {
		log.Printf("cargar página principal: %v", err)
		http.Error(w, "No se pudieron cargar los datos", http.StatusInternalServerError)
		return
	}
	app.renderIndex(w, r, data, http.StatusOK)
}

func (app *application) buildPageData(ctx context.Context) (views.PageData, error) {
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

	ventas, err := app.consultas.ListVentasConDetalle(ctx)
	if err != nil {
		return views.PageData{}, fmt.Errorf("listar ventas: %w", err)
	}

	data := views.PageData{}
	ventasPorPublicacion := make(map[int64]db.ListVentasConDetalleRow, len(ventas))

	for _, venta := range ventas {
		ventasPorPublicacion[venta.IDPublicacion] = venta
		data.Ventas = append(data.Ventas, views.VentaVista{
			ID:            venta.VentaID,
			PublicacionID: venta.IDPublicacion,
			Comprador:     venta.CompradorNombre,
			Vendedor:      venta.VendedorNombre,
			Precio:        views.FormatearPrecio(venta.Precio),
			Fecha:         views.FormatearFecha(venta.Fecha),
		})
	}
	usuariosPorID := make(map[int64]db.Usuario, len(usuarios))

	for _, usuario := range usuarios {
		usuariosPorID[usuario.ID] = usuario
		nombre := views.NombreCompleto(usuario.Nombre, usuario.Apellido)
		data.Usuarios = append(data.Usuarios, views.UsuarioVista{
			ID:             usuario.ID,
			NombreCompleto: nombre,
			Contacto:       views.FormatearContacto(usuario.Contacto),
		})
		data.OpcionesUsuarios = append(data.OpcionesUsuarios, views.Opcion{
			ID:    usuario.ID,
			Texto: nombre,
		})
	}

	animalesPorID := make(map[int64]db.Animal, len(animales))

	for _, animal := range animales {
		animalesPorID[animal.ID] = animal
		dueno := fmt.Sprintf("Usuario #%d", animal.IDDueno)
		if usuario, ok := usuariosPorID[animal.IDDueno]; ok {
			dueno = views.NombreCompleto(usuario.Nombre, usuario.Apellido)
		}
		data.Animales = append(data.Animales, views.AnimalVista{
			ID:     animal.ID,
			Nombre: animal.Nombre,
			Raza:   animal.Raza,
			Dueno:  dueno,
			Precio: views.FormatearPrecio(animal.Precio),
		})
		data.OpcionesAnimales = append(data.OpcionesAnimales, views.Opcion{
			ID:    animal.ID,
			Texto: fmt.Sprintf("%s (%s)", animal.Nombre, animal.Raza),
		})
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

		comprador := ""
		fechaVenta := ""
		if venta, ok := ventasPorPublicacion[publicacion.ID]; ok {
			comprador = venta.CompradorNombre
			fechaVenta = views.FormatearFecha(venta.Fecha)
		}

		data.Publicaciones = append(data.Publicaciones, views.PublicacionVista{
			ID:         publicacion.ID,
			IDVendedor: publicacion.IDVendedor,
			Animal:     animalNombre,
			Raza:       raza,
			Vendedor:   vendedor,
			Precio:     views.FormatearPrecio(publicacion.Precio),
			Activa:     publicacion.Activo,
			Comprador:  comprador,
			FechaVenta: fechaVenta,
		})
	}

	return data, nil
}

func (app *application) renderIndex(w http.ResponseWriter, r *http.Request, data views.PageData, status int) {
	var buffer bytes.Buffer

	if err := views.IndexPage(data).Render(r.Context(), &buffer); err != nil {
		log.Printf("renderizar página principal: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Printf("escribir respuesta: %v", err)
	}
}

func (app *application) renderPageError(w http.ResponseWriter, r *http.Request, mensaje string, status int) {
	data, err := app.buildPageData(r.Context())
	if err != nil {
		log.Printf("cargar datos para error: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	data.Error = mensaje
	app.renderIndex(w, r, data, status)
}
