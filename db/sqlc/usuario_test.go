package db

import (
	"context"
	"database/sql"
	"testing"
)

// Con esta hacemos el crear el usuario
func TestCreateUsuario(t *testing.T) {
	ctx := context.Background()

	// parametros de entrada para el test:
	arg := CreateUsuarioParams{
		Nombre:   "Elsa",
		Apellido: "Popepe",
		Contacto: sql.NullString{
			String: "correo@example.com",
			Valid:  true,
		},
	}

	// Creamos el usuario
	usuario, err := testQueries.CreateUsuario(ctx, arg)

	// Chequemos que el usuario este creado correctamente
	if err != nil {
		t.Fatalf("Fallo al crear Usuario: %v", err)
	}

	// LIMPIEZA
	t.Cleanup(func() {
		_, err := testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", usuario.ID)
		if err != nil {
			t.Logf("Fallo al limpiar Usuario con ID %d: %v", usuario.ID, err)
		}
	})

	// Test Assert
	if usuario.ID == 0 {
		t.Errorf("Se esperaba un ID mayor a 0")
	}
	if usuario.Nombre != arg.Nombre {
		t.Errorf("Se esperaba: Nombre %s, se obtuvo %s", arg.Nombre, usuario.Nombre)
	}
	if usuario.Apellido != arg.Apellido {
		t.Errorf("Se esperaba: Apellido %s, se obtuvo %s", arg.Apellido, usuario.Apellido)
	}
	if usuario.Contacto.String != arg.Contacto.String {
		t.Errorf("Se esperaba: Contacto %s, se obtuvo %s", arg.Contacto.String, usuario.Contacto.String)
	}
}

// Con esta funcion chequeamos el get
func TestGetUsuario(t *testing.T) {
	ctx := context.Background()

	createArg := CreateUsuarioParams{
		Nombre:   "Jorge",
		Apellido: "Suspenso",
		Contacto: sql.NullString{
			String: "saaran123@example.com",
			Valid:  true,
		},
	}

	nuevoUsuario, err := testQueries.CreateUsuario(ctx, createArg)
	if err != nil {
		t.Fatalf("Fallo al crear el ussuario inicial: %v", err)
	}

	//lo que hace esta funcion es asegurarse de que cuando termine el test, borre lo que esta creando, incluso si es test falla, lo borra igual
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", nuevoUsuario.ID)
	})

	fetchedUsuario, err := testQueries.GetUsuario(ctx, nuevoUsuario.ID)

	// Assert
	if err != nil {
		t.Fatalf("Fallo al obtener el usuario con el ID %d: %v", nuevoUsuario.ID, err)
	}

	if fetchedUsuario.ID != nuevoUsuario.ID {
		t.Errorf("Esperado: ID %d, Obtenido: %d", nuevoUsuario.ID, fetchedUsuario.ID)
	}
	if fetchedUsuario.Nombre != nuevoUsuario.Nombre {
		t.Errorf("Esperado: Nombre %s, Obtenido: %s", nuevoUsuario.Nombre, fetchedUsuario.Nombre)
	}
}

// Con esta funcion actualizacion un usuario
func TestUpdateUsuario(t *testing.T) {
	ctx := context.Background()

	//usuario base
	createArg := CreateUsuarioParams{
		Nombre:   "Luis",
		Apellido: "Spinetta",
		Contacto: sql.NullString{String: "email-previo@exa.com", Valid: true},
	}
	usuarioBase, err := testQueries.CreateUsuario(ctx, createArg)
	if err != nil {
		t.Fatalf("Fallo al crear el usuario base: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", usuarioBase.ID)
	})

	// preparamos los datos actualizados y ejecutamos Update
	updateArg := UpdateUsuarioParams{
		ID:       usuarioBase.ID,
		Nombre:   "Luis",            //mantenemos nombre
		Apellido: "Almirante Brown", //cambiamos apellido
		Contacto: sql.NullString{String: "email-nuevo@artaud.com", Valid: true},
	}

	usuarioActualizado, err := testQueries.UpdateUsuario(ctx, updateArg)

	// 3. Assert
	if err != nil {
		t.Fatalf("Fallo al updatear Usuario: %v", err)
	}
	if usuarioActualizado.Nombre != createArg.Nombre {
		t.Errorf("Expected Nombre %s, got %s", createArg.Nombre, usuarioActualizado.Nombre) //se espera el nombre original
	}
	if usuarioActualizado.Apellido != updateArg.Apellido {
		t.Errorf("Expected Apellido %s, got %s", updateArg.Apellido, usuarioActualizado.Apellido) //se espera el apellido actualizado
	}
	if usuarioActualizado.Contacto.String != updateArg.Contacto.String {
		t.Errorf("Expected Contacto %s, got %s", updateArg.Contacto.String, usuarioActualizado.Contacto.String)
	}
}

// con esta funcion eliminamos un usuario
func TestDeleteUsuario(t *testing.T) {
	ctx := context.Background()

	//mock inicial
	createArg := CreateUsuarioParams{
		Nombre:   "Usuario",
		Apellido: "Funable",
		Contacto: sql.NullString{Valid: false},
	}
	usuarioABorrar, err := testQueries.CreateUsuario(ctx, createArg)
	if err != nil {
		t.Fatalf("Fallo al crear el usuario base : %v", err)
	}

	// delete
	err = testQueries.DeleteUsuario(ctx, usuarioABorrar.ID)
	if err != nil {
		t.Fatalf("Fallo al borrar Usuario: %v", err)
	}

	// intentamos buscarlo, DEBE fallar
	_, err = testQueries.GetUsuario(ctx, usuarioABorrar.ID)

	if err == nil {
		t.Error("Se esperaba un error, ya que el usuario fue eliminado, pero se obtuvo nil")
	}
	// error estándar de Go cuando un select no devuelve nada
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows, se obtuvo: %v", err)
	}
}
