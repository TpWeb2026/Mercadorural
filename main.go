package main

import (
	"fmt"
	"net/http" // es lo que me permite crear un servidor web y manejar peticiones HTTP
)

func main() {
	//Defino la direccion estatica
	// con el FileServer ya se contemplan las rutas inexistentes por lo cual al ingresar a una ruta que no existe devuelve 404 page not found
	// FileServer tambien contempla los Content-Type
	// creo un manejador de archivos del sistema de archivo de "/static"
	// el http.dir(staticdir) convierte la ruta del directorio en un sistema de archivo http
	staticdir := "./static"
	fs := http.FileServer(http.Dir(staticdir))

	// Registramos que se muestre la pagina cada vez que se accede ala raiz
	http.Handle("/", fs)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// Inicia el servidor HTTP
	err := http.ListenAndServe(port, nil) // sirve para arrancar un servidor web en el puerto especificado y manejar las peticiones entrantes. El segundo parámetro es nil porque no estoy usando un mux personalizado, sino el predeterminado.
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
