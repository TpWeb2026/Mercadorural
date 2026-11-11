package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	db "tp2/db/sqlc"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env; se usan las variables del entorno")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	conexion, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer conexion.Close()

	if err := conexion.Ping(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	log.Println("Conexión a la base de datos exitosa")

	app := &application{consultas: db.New(conexion), conexion: conexion}

	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", app.routes()); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
