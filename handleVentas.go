package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http" //este es para poder sacar el id
	"strconv"
	"strings"
	db "tp2/db/sqlc"
)

// esta estrucuta se hizo para poder comunicar y poder usar los metodos que nos creo sqlc
type estructuraVentas struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}

func (h *estructuraVentas) CRUDventas(w http.ResponseWriter, r *http.Request) {
	cantidadURL := strings.Split(r.URL.Path, "/")

	if len(cantidadURL) > 3 {
		http.Error(w, "URL no permitida", http.StatusBadRequest)
		return
	}

	// CASO CON ID: /Ventas/5
	if len(cantidadURL) == 3 && cantidadURL[2] != "" {
		switch r.Method {
		case http.MethodGet:
			h.getVentasId(w, r)
		case http.MethodPut:
			//h.actualizarVentasId(w, r) // DISCUTIR SI TIENE SENTIDO, YA QUE VENTAS ES UN HISTORICO Y NO TENDRIA SENTIDO ACTUALIZAR ALGO QUE ES PARA MANTENER UN REGISTRO
		case http.MethodDelete:
			//h.eliminarVentasId(w, r) // ESTE NO TIENE SENTIDO, YA QUE PARA QUE VAS A ELIMINAR UNA VENTA DE HISTORICO SI PUEDE ARRUINAR TU CONTABILIDAD
		default:
			http.Error(w, "Método no aceptado", http.StatusMethodNotAllowed)
		}
		return
	}

	// CASO SIN ID: /Ventas
	switch r.Method {
	case http.MethodGet:
		h.listaDeVentas(w, r)
	case http.MethodPost:
		h.agregarVentas(w, r)
	default:
		http.Error(w, "Método no aceptado", http.StatusMethodNotAllowed)
	}
}

func (h *estructuraVentas) listaDeVentas(w http.ResponseWriter, r *http.Request) {
	//guardamos en una variables todas las Ventas activas que nos devuelve la base de datos
	Ventas, err := h.Consultas.ListVentasConDetalle(r.Context())

	//chequeamos errores
	if err != nil {
		http.Error(w, "Error al obtener Ventas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// si no hay ninguna Ventas, entonces creamos una vacia para poder devolvesela al usuario en formato json
	if Ventas == nil {
		Ventas = []db.ListVentasConDetalleRow{} // ACA SE TENDRIA QUE DEVOLVER []DB.VENTUM{} PERO EL SQLC GENERO UNA ESTRUCTURA AL PEDIR LAS VENTAS ENTONCES TENEMOS QUE ADAPTARLO A ESA ESTRUCTURA
	}

	// se la devolvemos al usuario en formato json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Ventas)

}

func (h *estructuraVentas) getVentasId(w http.ResponseWriter, r *http.Request) {
	//sacamos el url que nos viene desde el request
	partesURL := strings.Split(r.URL.Path, "/")

	//lo convertimos a int64 porque esta en string
	id, err := strconv.ParseInt(partesURL[2], 10, 64)

	//chequeamos que el id sea correcto
	if err != nil {
		http.Error(w, "ID inválido en la URL", http.StatusBadRequest)
		return
	}

	//llamamos a la funcion que nos devuelve la ventas dependiendo de un id
	venta, err := h.Consultas.GetVenta(r.Context(), id)

	//chequeamos los errores
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No se encontró la venta con ese ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//encapsulamos la respuesta en formato json y se la mostramos al usuario
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(venta)
}

func (h *estructuraVentas) agregarVentas(w http.ResponseWriter, r *http.Request) {
	var Venta db.Ventum
	//aca lo que hacemos es decodificar lo que nos mando el cliente por la solicitud
	err := json.NewDecoder(r.Body).Decode(&Venta)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ACA TENEMOS QUE VALIDAR QUE LOS ATRIBUTOS DEL ANIMAL SEAN VALIDOS, ESTO SERIA LAS REGLAS DE NEGOCIO
	if Venta.IDComprador < 0 || Venta.IDVendedor < 0 || Venta.IDPublicacion < 0 {
		http.Error(w, "El campo id comprador/vendedor/publicacion alguno no es valido", http.StatusBadRequest)
		return
	}

	//Validamos que el ID vendedor exista
	if _, err := h.Consultas.GetUsuario(r.Context(), Venta.IDVendedor); err != nil {
		http.Error(w, "El ID vendedor no existe", http.StatusBadRequest)
		return
	}
	//Validamos que el ID Comprador exista
	if _, err := h.Consultas.GetUsuario(r.Context(), Venta.IDComprador); err != nil {
		http.Error(w, "El ID comprador no existe", http.StatusBadRequest)
		return
	}
	//Validamos que el ID de publicacion exista
	if _, err := h.Consultas.GetPublicacion(r.Context(), Venta.IDPublicacion); err != nil {
		http.Error(w, "El ID de publicacion no existe", http.StatusBadRequest)
		return
	}

	// mapeamos los datos a los parametros que genero sqlc
	ventaBase := db.CreateVentaParams{
		IDPublicacion: Venta.IDPublicacion,
		IDVendedor:    Venta.IDVendedor,
		IDComprador:   Venta.IDComprador,
	}

	// aca lo que hacemos es ejecutar la insercion de la venta nueva a la base de datos
	ventaNueva, err := h.Consultas.CreateVenta(r.Context(), ventaBase)

	//chequemos para ver si se pudo completar la insercion o no
	if err != nil {
		http.Error(w, "No se pudo guardar la venta en la base de datos"+err.Error(), http.StatusInternalServerError)
		return
	}

	//ahora le respondemos al cliente con el http de animal creado y el nuevo objeto
	w.Header().Set("Content-Type", "application/json")

	//le notificamos que se creo bien con el codigo 201
	w.WriteHeader(http.StatusCreated)

	// lo devolvemos en formato json
	json.NewEncoder(w).Encode(ventaNueva)
}
