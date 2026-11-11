package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Driver
)

// Declaramos estas variables de forma global para que todos los archivos _test.go que creemos puedan usarlas.
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// conexion segun docker-compose.yml
	dsn := "postgresql://postgres:secretoclave@localhost:5432/tp2?sslmode=disable"

	// Abrimos la conexión
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos de test: %v", err)
	}

	// Inicializamos la estructura de sqlc
	testQueries = New(testDB)

	//lanza la ejecución de todos los tests
	code := m.Run()

	// cuando todos los tests terminan, cerramos la conexión
	testDB.Close()
	os.Exit(code)
}
