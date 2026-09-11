package db

import (
	"context"
	"database/sql"
	"testing"
)

// Función auxiliar interna para crear el Usuario y Animal requeridos por las Claves Foráneas
// aca esta declarado que la funcion crear usuario y animal te devuelve dos cosas, un usuario y un animal
func crearUsuarioYAnimalPrevios(t *testing.T, ctx context.Context) (Usuario, Animal) {
	// Primero, creamos un usuario y despues un animal, para poder crear la publicacion, porque necesita de un animal y de un usuario

	//usamos el structu para crear el usuario
	creacionUsuario := CreateUsuarioParams{
		Nombre:   "Carlos",
		Apellido: "Vendedor",
		Contacto: sql.NullString{String: "carlos@test.com", Valid: true},
	}

	//chequeamos que ele usuario se haya creado
	usuario, err := testQueries.CreateUsuario(ctx, creacionUsuario)
	if err != nil {
		t.Fatalf("Fallo al crear usuario previo para la FK: %v", err)
	}

	// 2. Registrar limpieza de Usuario (se ejecutará al final del test)
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", usuario.ID)
	})

	//Creamos el animal

	animalNuevo := CreateAnimalParams{
		Nombre:  "Rex",
		Raza:    sql.NullString{String: "Pastor Aleman", Valid: true},
		IDDueno: usuario.ID,
		Precio:  "45000.00",
	}

	animal, err := testQueries.CreateAnimal(ctx, animalNuevo)

	if err != nil {
		t.Fatalf("Fallo al crear animal previo para la FK: %v", err)
	}

	// 4. Registrar limpieza de Animal (se ejecutará ANTES que la del Usuario por el orden LIFO de Cleanup)
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Animal WHERE id = $1", animal.ID)
	})

	return usuario, animal
}

// funcion de testing para la creacion de una publicacion
func TestCreatePublicacion(t *testing.T) {
	ctx := context.Background()
	//asignamos lo que nos devuelve la funcion a una variable de usuario y a una variable de animal
	usuario, animal := crearUsuarioYAnimalPrevios(t, ctx)

	arg := CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     "45000.00",
		IDVendedor: usuario.ID,
	}

	pub, err := testQueries.CreatePublicacion(ctx, arg)
	if err != nil {
		t.Fatalf("Fallo al crear Publicacion: %v", err)
	}

	// Limpieza de la publicación
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Publicacion WHERE id = $1", pub.ID)
	})

	// Aca hacemos el assert, para asegurarnos que todo este bien, antes de seguir con los otros metodos
	if pub.ID == 0 {
		t.Errorf("Se esperaba un ID mayor a 0 para la publicacion")
	}
	if pub.IDAnimal != arg.IDAnimal {
		t.Errorf("Esperado IdAnimal %d, se obtuvo %d", arg.IDAnimal, pub.IDAnimal)
	}
	if pub.IDVendedor != arg.IDVendedor {
		t.Errorf("Esperado IdVendedor %d, se obtuvo %d", arg.IDVendedor, pub.IDVendedor)
	}
	if pub.Precio != arg.Precio {
		t.Errorf("Esperado Precio %s, se obtuvo %s", arg.Precio, pub.Precio)
	}
}

// funcion de testing para hacerle el get a la publicacion
func TestGetPublicacion(t *testing.T) {
	ctx := context.Background()
	//volvemos a referenciar el usuario y el animal creado en la funcion de crearUsuarioYanimal
	usuario, animal := crearUsuarioYAnimalPrevios(t, ctx)

	publi := CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     "50000.00",
		IDVendedor: usuario.ID,
	}

	pubCreada, err := testQueries.CreatePublicacion(ctx, publi)
	if err != nil {
		t.Fatalf("Fallo al crear publicacion base: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Publicacion WHERE id = $1", pubCreada.ID)
	})

	getPubli, err := testQueries.GetPublicacion(ctx, pubCreada.ID)
	if err != nil {
		t.Fatalf("Fallo al obtener la publicacion con ID %d: %v", pubCreada.ID, err)
	}

	if getPubli.ID != pubCreada.ID {
		t.Errorf("Esperado ID %d, Obtenido %d", pubCreada.ID, getPubli.ID)
	}
}

// funcion de testing para actualizar publicacion
func TestUpdatePublicacion(t *testing.T) {
	ctx := context.Background()
	usuario, animal := crearUsuarioYAnimalPrevios(t, ctx)

	crearPubli := CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     "50000.00",
		IDVendedor: usuario.ID,
	}
	pubBase, err := testQueries.CreatePublicacion(ctx, crearPubli)

	if err != nil {
		t.Fatalf("Fallo al crear publicacion base: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Publicacion WHERE id = $1", pubBase.ID)
	})

	actualizarPublicacion := UpdatePublicacionParams{
		ID:         pubBase.ID,
		IDAnimal:   animal.ID,
		Precio:     "40000.00", // Precio rebajado
		IDVendedor: usuario.ID,
	}

	pubActualizada, err := testQueries.UpdatePublicacion(ctx, actualizarPublicacion)
	if err != nil {
		t.Fatalf("Fallo al actualizar publicacion: %v", err)
	}

	if pubActualizada.Precio != actualizarPublicacion.Precio {
		t.Errorf("Esperado precio %s, se obtuvo %s", actualizarPublicacion.Precio, pubActualizada.Precio)
	}
}

//funcion testing para eliminar una publicacion

func TestDeletePublicacion(t *testing.T) {
	ctx := context.Background()
	usuario, animal := crearUsuarioYAnimalPrevios(t, ctx)
	publicacion := CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     "30000.00",
		IDVendedor: usuario.ID,
	}

	pub, err := testQueries.CreatePublicacion(ctx, publicacion)
	if err != nil {
		t.Fatalf("Fallo al crear la publicacion a eliminar: %v", err)
	}

	// Ejecutamos Delete
	err = testQueries.DeletePublicacion(ctx, pub.ID)
	if err != nil {
		t.Fatalf("Fallo al borrar Publicacion: %v", err)
	}

	// Comprobamos que el Get devuelva sql.ErrNoRows, para asegurar que se elimino correctamente
	_, err = testQueries.GetPublicacion(ctx, pub.ID)
	if err == nil {
		t.Error("Se esperaba un error ya que la publicacion fue eliminada, pero se obtuvo nil")
	}
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows, se obtuvo: %v", err)
	}
}
