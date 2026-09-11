package test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	// Importación usando el módulo "tp2" definido en tu go.mod
	db "tp2/db/sqlc"
)

// esta constante es para conectarse con el contenedor
const conexion = "postgres://postgres:secretoclave@localhost:5432/tp2?sslmode=disable"

func TestUsuario_CRUD(t *testing.T) {
	//intento la conexion con postgresql
	conn, err := sql.Open("postgres", conexion)
	if err != nil {
		t.Fatalf("Error de conexion a la base de datos: %v", err)
	}

	defer conn.Close()

	// con la siguiente instancia, asocia la conexion de la base de datos con las consultas sqlc
	queries := db.New(conn)
	ctx := context.Background()

	// 1. Primer paso, creacion del usuario, para chequear que ande bien el codigo que me hizo sqlc generate, y que se puede usar
	usuarioCreado, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Nombre:   "Pepito",
		Apellido: "Sope",
		Contacto: sql.NullString{String: "pepito@exa.edu.com", Valid: true},
		// se usa sql.NullString, porque el contacto en la base de datos puede ser nulo, se hace eso
	})
	if err != nil {
		t.Fatalf("Error al crear usuario: %v", err)
	}

	// 2. Segundo paso, chequear que el usuario se haya creado correctamente
	usuarioObtenido, err := queries.GetUsuario(ctx, usuarioCreado.ID)
	if err != nil {
		t.Fatalf("Error al obtener usuario: %v", err)
	}
	if usuarioObtenido.Nombre != "Pepito" {
		t.Errorf("El nombre no coincide. Esperado: Pepito, Obtenido: %s", usuarioObtenido.Nombre)
	}

	// 3. Tercer paso, actualizacion de informacion
	_, err = queries.UpdateUsuario(ctx, db.UpdateUsuarioParams{
		ID:       usuarioCreado.ID,
		Nombre:   "Ernesto Julian",
		Apellido: "Gomez",
		Contacto: sql.NullString{String: "Ernestopiola@test.com", Valid: true},
	})
	if err != nil {
		t.Fatalf("Error al actualizar usuario: %v", err)
	}

	// 4. Cuarto paso, eliminacion de usuario
	err = queries.DeleteUsuario(ctx, usuarioCreado.ID)
	if err != nil {
		t.Fatalf("Error al eliminar usuario: %v", err)
	}
}
