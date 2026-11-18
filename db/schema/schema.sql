CREATE TABLE Usuario (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    apellido VARCHAR(50) NOT NULL,
    contacto VARCHAR(150),
    activo BOOLEAN NOT NULL DEFAULT true    /*Campo añadido para permitir el soft delete y mantener
                                             las transacciones historicas sin sobrecargar la tabla venta*/
);

CREATE TABLE Animal (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    raza VARCHAR(50) NOT NULL,  /*ACA LO AGREGAMOS PARA EL TP3 PORQUE NOS DIMOS CUENTA QUE NO SE PUEDE TENER ANIMAL SIN RAZA*/
    id_dueno BIGINT NOT NULL,
    precio NUMERIC(10, 2) NOT NULL,

    CONSTRAINT fk_animal_dueno 
        FOREIGN KEY (id_dueno) REFERENCES Usuario(id) 
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Publicacion (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id_animal BIGINT NOT NULL,
    precio NUMERIC(10, 2) NOT NULL,
    id_vendedor BIGINT NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT fk_pub_animal 
        FOREIGN KEY (id_animal) REFERENCES Animal(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_pub_vendedor 
        FOREIGN KEY (id_vendedor) REFERENCES Usuario(id) 
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Venta (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id_publicacion BIGINT NOT NULL,
    id_vendedor BIGINT NOT NULL,
    id_comprador BIGINT NOT NULL,
    fecha TIMESTAMPTZ DEFAULT NOW(),

    /*  Ver el tema de cuando se elimina un usuario, la publicacion tendria que quedarse, lo mismo con el animal */
    /* sacarle las restrcciones y dejarla solo, sin fk, para no tener problema con el tema de eliminar los usuarios */
    CONSTRAINT fk_venta_pub 
        FOREIGN KEY (id_publicacion) REFERENCES Publicacion(id) 
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_venta_vendedor 
        FOREIGN KEY (id_vendedor) REFERENCES Usuario(id) 
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_venta_comprador 
        FOREIGN KEY (id_comprador) REFERENCES Usuario(id) 
        ON DELETE RESTRICT ON UPDATE CASCADE
);