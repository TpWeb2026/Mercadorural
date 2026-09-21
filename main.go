package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http" // es lo que me permite crear un servidor web y manejar peticiones HTTP
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //driver de postgres
)

func main() {
	//tp3 ejercicio 1, establecer conexion con la DB
	// godotenv.Load() lee tu archivo .env y mete esos datos ocultos en el sistema operativo
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontro el archivo .env")
	}
	/*// os.Getenv("DB_DSN") busca específicamente el dato guardado bajo la etiqueta "DB_DSN" en el .env
	credenciales := os.Getenv("DB_DSN")*/

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	if dsn == "" {
		log.Fatal("Error la variable DB_DSN no esta configurada")
	}
	//establecemos conexion con la DB
	conexion, errDB := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error al abrir la conexion", errDB)
	}
	errDB = conexion.Ping()
	if errDB != nil {
		log.Fatal("Error conectando ala base de datos", errDB)
	}
	fmt.Println("Conexion ala base de datos exitosa")

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
	port = ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// Inicia el servidor HTTP
	err = http.ListenAndServe(port, nil) // sirve para arrancar un servidor web en el puerto especificado y manejar las peticiones entrantes. El segundo parámetro es nil porque no estoy usando un mux personalizado, sino el predeterminado.
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
