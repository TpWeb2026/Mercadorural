package views

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Opcion struct {
	ID    int64
	Texto string
}

type UsuarioVista struct {
	ID             int64
	NombreCompleto string
	Contacto       string
}

type AnimalVista struct {
	ID     int64
	Nombre string
	Raza   string
	Dueno  string
	Precio string
}

type PublicacionVista struct {
	ID         int64
	IDVendedor int64
	Animal     string
	Raza       string
	Vendedor   string
	Precio     string
	Activa     bool
	Comprador  string
	FechaVenta string
}

type VentaVista struct {
	ID            int64
	PublicacionID int64
	Comprador     string
	Vendedor      string
	Precio        string
	Fecha         string
}

type UsuarioFormData struct {
	Nombre   string
	Apellido string
	Contacto string
}

type AnimalFormData struct {
	Nombre  string
	Raza    string
	IDDueno int64
	Precio  string
}

type PublicacionFormData struct {
	IDAnimal   int64
	IDVendedor int64
	Precio     string
}

type PageData struct {
	Usuarios         []UsuarioVista
	Animales         []AnimalVista
	Publicaciones    []PublicacionVista
	Ventas           []VentaVista
	OpcionesUsuarios []Opcion
	OpcionesAnimales []Opcion
	UsuarioForm      UsuarioFormData
	AnimalForm       AnimalFormData
	PublicacionForm  PublicacionFormData
	Error            string
}

func NombreCompleto(nombre, apellido string) string {
	return strings.TrimSpace(nombre + " " + apellido)
}

func FormatearContacto(contacto sql.NullString) string {
	if contacto.Valid && contacto.String != "" {
		return contacto.String
	}
	return "Sin contacto"
}

func FormatearPrecio(precio string) string {
	valor, err := strconv.ParseFloat(precio, 64)
	if err != nil {
		return precio
	}
	return fmt.Sprintf("$ %.2f", valor)
}

func FormatearFecha(fecha sql.NullTime) string {
	if !fecha.Valid {
		return "Sin fecha"
	}
	return fecha.Time.Format("02/01/2006 15:04")
}

func compradoresDisponibles(usuarios []Opcion, vendedorID int64) []Opcion {
	compradores := make([]Opcion, 0, len(usuarios))
	for _, usuario := range usuarios {
		if usuario.ID != vendedorID {
			compradores = append(compradores, usuario)
		}
	}
	return compradores
}

func urlEliminar(tipo string, id int64) string {
	return fmt.Sprintf("/%s/%d/eliminar", tipo, id)
}

func urlComprar(id int64) string {
	return fmt.Sprintf("/publicacion/%d/comprar", id)
}
