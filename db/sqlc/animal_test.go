package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestCreateAnimal(t *testing.T) {
	ctx := context.Background()

	dueno, err := testQueries.CreateUsuario(ctx, CreateUsuarioParams{
		Nombre:   "Dueño",
		Apellido: "Test",
	})
	if err != nil {
		t.Fatalf("Fallo al crear usuario: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Animal")
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", dueno.ID)
	})

	createArg := CreateAnimalParams{
		Nombre:  "Bartolito",
		Raza:    "mitre",
		IDDueno: dueno.ID,
		Precio:  "100.00",
	}

	nuevoAnimal, err := testQueries.CreateAnimal(ctx, createArg)
	if err != nil {
		t.Fatalf("Fallo al crear el animal inicial: %v", err)
	}

	if nuevoAnimal.ID == 0 {
		t.Errorf("Se esperaba un ID mayor a 0")
	}

	if nuevoAnimal.Nombre != createArg.Nombre {
		t.Errorf("Se esperaba: Nombre %s, se obtuvo %s", createArg.Nombre, nuevoAnimal.Nombre)
	}

	if nuevoAnimal.Raza != createArg.Raza {
		t.Errorf("Se esperaba: Raza %s, se obtuvo %s", createArg.Raza, nuevoAnimal.Raza)
	}

	if nuevoAnimal.IDDueno != createArg.IDDueno {
		t.Errorf("Se esperaba: ID Dueno %d, se obtuvo %d", createArg.IDDueno, nuevoAnimal.IDDueno)
	}

	if nuevoAnimal.Precio != createArg.Precio {
		t.Errorf("Se esperaba Precio: %s, se obtuvo %s", createArg.Precio, nuevoAnimal.Precio)
	}
}

func TestGetAnimal(t *testing.T) {
	ctx := context.Background()

	dueno, err := testQueries.CreateUsuario(ctx, CreateUsuarioParams{
		Nombre:   "Dueño",
		Apellido: "Test",
	})
	if err != nil {
		t.Fatalf("Fallo al crear usuario: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Animal")
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", dueno.ID)
	})

	createArg := CreateAnimalParams{
		Nombre:  "Bartolito",
		Raza:    "calos",
		IDDueno: dueno.ID,
		Precio:  "100.00",
	}

	nuevoAnimal, err := testQueries.CreateAnimal(ctx, createArg)
	if err != nil {
		t.Fatalf("Fallo al crear el animal inicial: %v", err)
	}

	fetchedAnimal, err := testQueries.GetAnimal(ctx, nuevoAnimal.ID)
	if err != nil {
		t.Fatalf("Fallo al obtener el animal ID %d: %v", nuevoAnimal.ID, err)
	}

	if fetchedAnimal.Nombre != nuevoAnimal.Nombre {
		t.Errorf("Se esperaba: Nombre %s, se obtuvo %s", nuevoAnimal.Nombre, fetchedAnimal.Nombre)
	}

	if fetchedAnimal.Raza != nuevoAnimal.Raza {
		t.Errorf("Se esperaba: Raza %s, se obtuvo %s", nuevoAnimal.Raza, fetchedAnimal.Raza)
	}

	if fetchedAnimal.IDDueno != nuevoAnimal.IDDueno {
		t.Errorf("Se esperaba: ID Dueno %d, se obtuvo %d", nuevoAnimal.IDDueno, fetchedAnimal.IDDueno)
	}

	if fetchedAnimal.Precio != nuevoAnimal.Precio {
		t.Errorf("Se esperaba Precio: %s, se obtuvo %s", nuevoAnimal.Precio, fetchedAnimal.Precio)
	}
}

func TestUpdateAnimal(t *testing.T) {
	ctx := context.Background()

	usuarioArg := CreateUsuarioParams{
		Nombre:   "Dueño",
		Apellido: "Update",
		Contacto: sql.NullString{Valid: false},
	}
	dueno, err := testQueries.CreateUsuario(ctx, usuarioArg)
	if err != nil {
		t.Fatalf("Fallo al inicializar el usuario: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Animal")
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", dueno.ID)
	})

	animalArg := CreateAnimalParams{
		Nombre:  "Firulais",
		Raza:    "alamadre",
		IDDueno: dueno.ID,
		Precio:  "1500.00",
	}
	animalBase, err := testQueries.CreateAnimal(ctx, animalArg)
	if err != nil {
		t.Fatalf("Fallo al inicializar el Animal: %v", err)
	}

	updateArg := UpdateAnimalParams{
		ID:      animalBase.ID,
		Nombre:  animalBase.Nombre,  // se mantiene el nombre
		Raza:    "elvi",             // no mas Sabrositos
		IDDueno: animalBase.IDDueno, // se mantiene el dueño
		Precio:  "2500.00",
	}

	animalActualizado, err := testQueries.UpdateAnimal(ctx, updateArg)

	// 3. Assert
	if err != nil {
		t.Fatalf("Fallo al actualizar el animal: %v", err)
	}
	if animalActualizado.Raza != updateArg.Raza {
		t.Errorf("Esperado: Raza %s, obtenido: %s", updateArg.Raza, animalActualizado.Raza)
	}
	if animalActualizado.Precio != updateArg.Precio {
		t.Errorf("Esperado: Precio %s, obtenido: %s", updateArg.Precio, animalActualizado.Precio)
	}
}

func TestDeleteAnimal(t *testing.T) {
	ctx := context.Background()

	dueno, _ := testQueries.CreateUsuario(ctx, CreateUsuarioParams{
		Nombre:   "Dueño",
		Apellido: "Delete",
	})

	animalABorrar, err := testQueries.CreateAnimal(ctx, CreateAnimalParams{
		Nombre:  "Borrable",
		IDDueno: dueno.ID,
		Precio:  "100.00",
	})
	if err != nil {
		t.Fatalf("Fallo al crear: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", dueno.ID)
	})

	err = testQueries.DeleteAnimal(ctx, animalABorrar.ID)
	if err != nil {
		t.Fatalf("Fallo al borrar animal: %v", err)
	}

	_, err = testQueries.GetAnimal(ctx, animalABorrar.ID)
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba error, se obtuvo:  %v", err)
	}
}
