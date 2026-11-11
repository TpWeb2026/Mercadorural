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

type OperadorVista struct {
	ID             int64
	NombreCompleto string
	Contacto       string
}

type PublicacionVista struct {
	ID       int64
	Animal   string
	Raza     string
	Vendedor string
	Precio   string
}

type UsuarioFormData struct {
	Nombre   string
	Apellido string
	Contacto string
}

type PublicacionFormData struct {
	NombreAnimal string
	Raza         string
	Precio       string
}

type EdicionPublicacionData struct {
	PublicacionID int64
	IDAnimal      int64
	NombreAnimal  string
	Raza          string
	Precio        string
}

type PageData struct {
	Operador             *OperadorVista
	Usuarios             []Opcion
	PublicacionesPropias []PublicacionVista
	PublicacionesAjenas  []PublicacionVista
	PublicacionForm      PublicacionFormData
	UsuarioForm          UsuarioFormData
	Error                string
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

func urlEliminar(tipo string, id int64) string {
	return fmt.Sprintf("/%s/%d/eliminar", tipo, id)
}

func urlComprar(id int64) string {
	return fmt.Sprintf("/publicacion/%d/comprar", id)
}

func urlEditar(id int64) string {
	return fmt.Sprintf("/publicacion/%d/editar", id)
}
