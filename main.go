package main

import (
	"fmt"
	"net/http" // es lo que me permite crear un servidor web y manejar peticiones HTTP
)

func Mostrapag(w http.ResponseWriter, r *http.Request) {
	// aca controlo la ruta inexistentes para que me de error, sino con el serverfile por mas que pongas una ruta inexistente te da lo del index.html
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// si bien serverfile ya maneja internamente los Content-Type osea no es necesario hacerlo manual, pero lo dejo escrito por si acaso se pedia eso
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// el serverfile sirve para un archivo, en cambio el fileserver es para server a archivos estaticos
	http.ServeFile(w, r, "index.html")
}

func main() {
	// Registramos que se muestre la pagina cada vez que se accede ala raiz
	http.HandleFunc("/", Mostrapag)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// Inicia el servidor HTTP
	err := http.ListenAndServe(port, nil) // sirve para arrancar un servidor web en el puerto especificado y manejar las peticiones entrantes. El segundo parámetro es nil porque no estoy usando un mux personalizado, sino el predeterminado.
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
