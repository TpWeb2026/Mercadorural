package db

import (
	"context"
	"database/sql"
	"testing"
)

// Función auxiliar para preparar toda la cadena de dependencias (FKs)

func crearEntidadesPreviasParaVenta(t *testing.T, ctx context.Context) (Usuario, Usuario, Publicacion) {

	// Con esta funcion creamos un vendedor
	vende := CreateUsuarioParams{
		Nombre:   "Marito",
		Apellido: "Carlitos",
		Contacto: sql.NullString{String: "mario@vende.com", Valid: true},
	}

	vendedor, err := testQueries.CreateUsuario(ctx, vende)

	if err != nil {
		t.Fatalf("Fallo al crear usuario vendedor: %v", err)
	}
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", vendedor.ID)
	})

	//Creamos el comprador
	compra := CreateUsuarioParams{
		Nombre:   "Lucia",
		Apellido: "Estela",
		Contacto: sql.NullString{String: "lucia@compradora.com", Valid: true},
	}

	// Con esta funcion crear Comprador
	comprador, err := testQueries.CreateUsuario(ctx, compra)
	if err != nil {
		t.Fatalf("Fallo al crear usuario comprador: %v", err)
	}
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", comprador.ID)
	})

	// Creamos un Animal asociado al Vendedor
	creamosAnimal := CreateAnimalParams{
		Nombre:  "Tormenta",
		Raza:    sql.NullString{String: "Caballo Criollo", Valid: true},
		IDDueno: vendedor.ID, // aca sacamos el id que tiene el vendedor que creamos anteriomente
		Precio:  "120000.00",
	}

	animal, err := testQueries.CreateAnimal(ctx, creamosAnimal)
	if err != nil {
		t.Fatalf("Fallo al crear animal previo: %v", err)
	}
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Animal WHERE id = $1", animal.ID)
	})

	// Creamos una publicacion asociada al animal que creamos anteriormente y al vendedor que tambien creamos anteriormete

	creamosPublicacion := CreatePublicacionParams{
		IDAnimal:   animal.ID,
		Precio:     "120000.00",
		IDVendedor: vendedor.ID,
	}

	publicacion, err := testQueries.CreatePublicacion(ctx, creamosPublicacion)
	if err != nil {
		t.Fatalf("Fallo al crear publicacion previa: %v", err)
	}
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Publicacion WHERE id = $1", publicacion.ID)
	})

	return vendedor, comprador, publicacion //aca retornamos el vendedor, el comprador y la publicacion que creamos para poder hacer los test correspondientes
}

// Con esta funcion, creamos la venta
func TestCreateVenta(t *testing.T) {
	ctx := context.Background()
	vendedor, comprador, publicacion := crearEntidadesPreviasParaVenta(t, ctx) //creamos variables para que tomen los valores que nos devuelve la funcion donde creamos el usuaio, la publicacion y la venta

	nuevaVenta := CreateVentaParams{
		IDPublicacion: publicacion.ID,
		IDVendedor:    vendedor.ID,
		IDComprador:   comprador.ID,
	}

	venta, err := testQueries.CreateVenta(ctx, nuevaVenta)
	if err != nil {
		t.Fatalf("Fallo al crear la Venta: %v", err)
	}

	// Limpieza de la Venta creada
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Venta WHERE id = $1", venta.ID)
	})

	// Asserts
	if venta.ID == 0 {
		t.Errorf("Se esperaba un ID mayor a 0 para la venta")
	}
	if venta.IDPublicacion != nuevaVenta.IDPublicacion {
		t.Errorf("Esperado IDPublicacion %d, se obtuvo %d", nuevaVenta.IDPublicacion, venta.IDPublicacion)
	}
	if venta.IDVendedor != nuevaVenta.IDVendedor {
		t.Errorf("Esperado IDVendedor %d, se obtuvo %d", nuevaVenta.IDVendedor, venta.IDVendedor)
	}
	if venta.IDComprador != nuevaVenta.IDComprador {
		t.Errorf("Esperado IDComprador %d, se obtuvo %d", nuevaVenta.IDComprador, venta.IDComprador)
	}
}

// Creamos la funcion que nos da la venta
func TestGetVenta(t *testing.T) {
	ctx := context.Background()
	vendedor, comprador, publicacion := crearEntidadesPreviasParaVenta(t, ctx)

	ventaCreada, err := testQueries.CreateVenta(ctx, CreateVentaParams{
		IDPublicacion: publicacion.ID,
		IDVendedor:    vendedor.ID,
		IDComprador:   comprador.ID,
	})
	if err != nil {
		t.Fatalf("Fallo al crear la venta inicial: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Venta WHERE id = $1", ventaCreada.ID)
	})

	// Chequeamos el get venta
	chequearVenta, err := testQueries.GetVenta(ctx, ventaCreada.ID)
	if err != nil {
		t.Fatalf("Fallo al obtener la venta con ID %d: %v", ventaCreada.ID, err)
	}

	//chequeamos que el get nos trae el id que acabamos de crear
	if chequearVenta.ID != ventaCreada.ID {
		t.Errorf("Esperado ID %d, Obtenido: %d", ventaCreada.ID, chequearVenta.ID)
	}

	// Hacemos el test de la funcion de listVentasByComprador
	ventasComprador, err := testQueries.ListVentasByComprador(ctx, comprador.ID)
	if err != nil || len(ventasComprador) == 0 {
		t.Fatalf("Fallo al listar ventas por comprador: %v", err)
	}

	// Test ListVentasConDetalle
	ventasDetalle, err := testQueries.ListVentasConDetalle(ctx)
	if err != nil || len(ventasDetalle) == 0 {
		t.Fatalf("Fallo al listar ventas con detalle: %v", err)
	}
}

// con esta funcion chequemos que ande bien el tema de actualizar alguna venta
func TestUpdateVenta(t *testing.T) {
	ctx := context.Background()
	vendedor, compradorOriginal, publicacion := crearEntidadesPreviasParaVenta(t, ctx)

	ventaBase, err := testQueries.CreateVenta(ctx, CreateVentaParams{
		IDPublicacion: publicacion.ID,
		IDVendedor:    vendedor.ID,
		IDComprador:   compradorOriginal.ID,
	})
	if err != nil {
		t.Fatalf("Fallo al crear venta base: %v", err)
	}

	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Venta WHERE id = $1", ventaBase.ID)
	})

	// Creamos un segundo comprador para transferir la venta en la prueba
	nuevoComprador, err := testQueries.CreateUsuario(ctx, CreateUsuarioParams{
		Nombre:   "Pedrito",
		Apellido: "tuvi",
		Contacto: sql.NullString{String: "pedro@test.com", Valid: true},
	})
	if err != nil {
		t.Fatalf("Fallo al crear nuevo comprador para update: %v", err)
	}
	t.Cleanup(func() {
		testDB.ExecContext(ctx, "DELETE FROM Usuario WHERE id = $1", nuevoComprador.ID)
	})

	updateArg := UpdateVentaParams{
		ID:            ventaBase.ID,
		IDPublicacion: publicacion.ID,
		IDVendedor:    vendedor.ID,
		IDComprador:   nuevoComprador.ID, // Actualizamos el comprador
	}

	ventaActualizada, err := testQueries.UpdateVenta(ctx, updateArg)
	if err != nil {
		t.Fatalf("Fallo al actualizar la venta: %v", err)
	}

	if ventaActualizada.IDComprador != nuevoComprador.ID {
		t.Errorf("Esperado IDComprador %d, se obtuvo %d", nuevoComprador.ID, ventaActualizada.IDComprador)
	}
}

// con esta funcion chequemos que ande bien el tema de eliminar una venta
func TestDeleteVenta(t *testing.T) {
	ctx := context.Background()
	vendedor, comprador, publicacion := crearEntidadesPreviasParaVenta(t, ctx)

	venta, err := testQueries.CreateVenta(ctx, CreateVentaParams{
		IDPublicacion: publicacion.ID,
		IDVendedor:    vendedor.ID,
		IDComprador:   comprador.ID,
	})
	if err != nil {
		t.Fatalf("Fallo al crear la venta a eliminar: %v", err)
	}

	// Ejecutamos Delete
	err = testQueries.DeleteVenta(ctx, venta.ID)
	if err != nil {
		t.Fatalf("Fallo al eliminar Venta: %v", err)
	}

	// Intentamos buscar la venta eliminada; debe devolver sql.ErrNoRows
	_, err = testQueries.GetVenta(ctx, venta.ID)
	if err == nil {
		t.Error("Se esperaba un error al buscar una venta eliminada, pero no hubo ninguno")
	}
	if err != sql.ErrNoRows {
		t.Errorf("Se esperaba sql.ErrNoRows, se obtuvo: %v", err)
	}
}
