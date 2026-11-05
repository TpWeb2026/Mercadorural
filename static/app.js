const API = {
  usuarios: "/usuario",
  animales: "/animales",
  publicaciones: "/publicacion",
};

const estado = {
  usuarios: [],
  animales: [],
  publicaciones: [],
};

const elementos = {};

const RECURSOS = {
  publicacion: {
    nombre: "la publicación",
    ruta: (id) => `${API.publicaciones}/${id}`,
    estado: () => elementos.publicacionStatus,
    exito: "Publicación eliminada correctamente.",
  },
  animal: {
    nombre: "el animal",
    ruta: (id) => `${API.animales}/${id}`,
    estado: () => elementos.animalStatus,
    exito: "Animal eliminado correctamente.",
  },
  usuario: {
    nombre: "el usuario",
    ruta: (id) => `${API.usuarios}/${id}`,
    estado: () => elementos.usuarioStatus,
    exito: "Usuario dado de baja correctamente.",
  },
};

function inicializarElementos() {
  elementos.appStatus = document.getElementById("app-status");

  elementos.formPublicacion = document.getElementById("form-publicacion");
  elementos.publicacionAnimal = document.getElementById("publicacion-animal");
  elementos.publicacionVendedor = document.getElementById("publicacion-vendedor");
  elementos.publicacionPrecio = document.getElementById("publicacion-precio");
  elementos.publicacionStatus = document.getElementById("publicacion-status");
  elementos.publicacionSubmit = document.getElementById("publicacion-submit");
  elementos.listaPublicaciones = document.getElementById("lista-publicaciones");
  elementos.contadorPublicaciones = document.getElementById("contador-publicaciones");

  elementos.formUsuario = document.getElementById("form-usuario");
  elementos.usuarioNombre = document.getElementById("usuario-nombre");
  elementos.usuarioApellido = document.getElementById("usuario-apellido");
  elementos.usuarioContacto = document.getElementById("usuario-contacto");
  elementos.usuarioStatus = document.getElementById("usuario-status");
  elementos.usuarioSubmit = document.getElementById("usuario-submit");
  elementos.listaUsuarios = document.getElementById("lista-usuarios");
  elementos.contadorUsuarios = document.getElementById("contador-usuarios");

  elementos.formAnimal = document.getElementById("form-animal");
  elementos.animalNombre = document.getElementById("animal-nombre");
  elementos.animalRaza = document.getElementById("animal-raza");
  elementos.animalDueno = document.getElementById("animal-dueno");
  elementos.animalPrecio = document.getElementById("animal-precio");
  elementos.animalStatus = document.getElementById("animal-status");
  elementos.animalSubmit = document.getElementById("animal-submit");
  elementos.listaAnimales = document.getElementById("lista-animales");
  elementos.contadorAnimales = document.getElementById("contador-animales");
}

async function apiRequest(ruta, opciones = {}) {
  const configuracion = {
    ...opciones,
    headers: { ...(opciones.headers || {}) },
  };

  if (configuracion.body !== undefined) {
    configuracion.headers["Content-Type"] = "application/json";
  }

  let respuesta;

  try {
    respuesta = await fetch(ruta, configuracion);
  } catch {
    throw new Error("No se pudo conectar con la API. Verificá que el servidor esté corriendo.");
  }

  if (respuesta.status === 204) {
    return null;
  }

  const texto = await respuesta.text();

  if (!respuesta.ok) {
    throw new Error(texto.trim() || `Error HTTP ${respuesta.status}`);
  }

  if (!texto) {
    return null;
  }

  try {
    return JSON.parse(texto);
  } catch {
    return texto;
  }
}

async function cargarTodo({ mostrarCarga = true } = {}) {
  if (mostrarCarga) {
    mostrarMensaje(elementos.listaPublicaciones, "Cargando publicaciones...");
    mostrarMensaje(elementos.listaUsuarios, "Cargando usuarios...", "empty-state empty-state--small");
    mostrarMensaje(elementos.listaAnimales, "Cargando animales...", "empty-state empty-state--small");
  }

  try {
    const [publicaciones, animales, usuarios] = await Promise.all([
      apiRequest(API.publicaciones),
      apiRequest(API.animales),
      apiRequest(API.usuarios),
    ]);

    estado.publicaciones = Array.isArray(publicaciones) ? publicaciones : [];
    estado.animales = Array.isArray(animales) ? animales : [];
    estado.usuarios = Array.isArray(usuarios) ? usuarios : [];

    limpiarStatus(elementos.appStatus);
    renderizarTodo();
  } catch (error) {
    mostrarStatus(elementos.appStatus, error.message, "error");
  }
}

function renderizarTodo() {
  renderizarSelects();
  renderizarPublicaciones();
  renderizarUsuarios();
  renderizarAnimales();
  actualizarContadores();
  actualizarDisponibilidadFormularios();
}

function renderizarSelects() {
  llenarSelect(
    elementos.publicacionAnimal,
    estado.animales,
    "Seleccioná un animal",
    (animal) => `${animal.nombre} (${animal.raza})`,
    "No hay animales cargados"
  );

  llenarSelect(
    elementos.publicacionVendedor,
    estado.usuarios,
    "Seleccioná un vendedor",
    (usuario) => nombreUsuario(usuario),
    "No hay usuarios cargados"
  );

  llenarSelect(
    elementos.animalDueno,
    estado.usuarios,
    "Seleccioná un dueño",
    (usuario) => nombreUsuario(usuario),
    "No hay usuarios cargados"
  );

  sincronizarVendedorConAnimal();
}

function llenarSelect(select, items, placeholder, etiqueta, textoVacio) {
  const valorPrevio = select.value;
  const fragmento = document.createDocumentFragment();

  if (!items.length) {
    const opcionVacia = document.createElement("option");
    opcionVacia.value = "";
    opcionVacia.textContent = textoVacio;
    opcionVacia.selected = true;
    fragmento.append(opcionVacia);
    select.replaceChildren(fragmento);
    select.disabled = true;
    return;
  }

  const opcionPlaceholder = document.createElement("option");
  opcionPlaceholder.value = "";
  opcionPlaceholder.textContent = placeholder;
  opcionPlaceholder.disabled = true;
  fragmento.append(opcionPlaceholder);

  for (const item of items) {
    const opcion = document.createElement("option");
    opcion.value = String(item.id);
    opcion.textContent = etiqueta(item);
    fragmento.append(opcion);
  }

  select.replaceChildren(fragmento);
  select.disabled = false;

  const existeValorPrevio = Array.from(select.options).some((opcion) => opcion.value === valorPrevio);
  select.value = existeValorPrevio ? valorPrevio : "";
}

function sincronizarVendedorConAnimal() {
  const animalId = Number(elementos.publicacionAnimal.value);

  if (!animalId) {
    return;
  }

  const animal = estado.animales.find((item) => item.id === animalId);

  if (!animal) {
    return;
  }

  const opcionDueno = Array.from(elementos.publicacionVendedor.options).find(
    (opcion) => Number(opcion.value) === animal.id_dueno
  );

  if (opcionDueno) {
    elementos.publicacionVendedor.value = opcionDueno.value;
  }
}

function renderizarPublicaciones() {
  if (!estado.publicaciones.length) {
    mostrarMensaje(
      elementos.listaPublicaciones,
      "No hay publicaciones activas. Creá la primera con el formulario."
    );
    return;
  }

  const fragmento = document.createDocumentFragment();

  for (const publicacion of estado.publicaciones) {
    fragmento.append(crearTarjetaPublicacion(publicacion));
  }

  elementos.listaPublicaciones.replaceChildren(fragmento);
}

function crearTarjetaPublicacion(publicacion) {
  const animal = estado.animales.find((item) => item.id === publicacion.id_animal);
  const vendedor = estado.usuarios.find((item) => item.id === publicacion.id_vendedor);

  const articulo = document.createElement("article");
  articulo.className = "publication-card";

  const cabecera = document.createElement("div");
  cabecera.className = "publication-card__top";

  const titulo = document.createElement("h3");
  titulo.textContent = animal ? animal.nombre : `Animal #${publicacion.id_animal}`;
  cabecera.append(titulo);

  const precio = document.createElement("span");
  precio.className = "price";
  precio.textContent = formatearPrecio(publicacion.precio);
  cabecera.append(precio);

  const detalles = document.createElement("ul");
  detalles.className = "publication-meta";
  detalles.append(crearDetalle(`Raza: ${animal ? animal.raza : "Sin datos"}`));
  detalles.append(
    crearDetalle(
      `Vendedor: ${vendedor ? nombreUsuario(vendedor) : `Usuario #${publicacion.id_vendedor}`}`
    )
  );
  detalles.append(crearDetalle(`Publicación #${publicacion.id}`));

  const acciones = document.createElement("div");
  acciones.className = "card-actions";
  acciones.append(crearBotonEliminar("publicacion", publicacion.id, "Eliminar publicación"));

  articulo.append(cabecera, detalles, acciones);
  return articulo;
}

function renderizarUsuarios() {
  if (!estado.usuarios.length) {
    mostrarMensaje(elementos.listaUsuarios, "Todavía no hay usuarios.", "empty-state empty-state--small");
    return;
  }

  const fragmento = document.createDocumentFragment();

  for (const usuario of estado.usuarios) {
    fragmento.append(crearItemUsuario(usuario));
  }

  elementos.listaUsuarios.replaceChildren(fragmento);
}

function crearItemUsuario(usuario) {
  const item = document.createElement("div");
  item.className = "mini-item";

  const info = document.createElement("div");
  info.className = "mini-item__info";

  const titulo = document.createElement("span");
  titulo.className = "mini-item__title";
  titulo.textContent = nombreUsuario(usuario);
  info.append(titulo);

  const detalle = document.createElement("span");
  detalle.className = "mini-item__detail";
  detalle.textContent = contactoUsuario(usuario);
  info.append(detalle);

  item.append(info, crearBotonEliminar("usuario", usuario.id, "Eliminar usuario"));
  return item;
}

function renderizarAnimales() {
  if (!estado.animales.length) {
    mostrarMensaje(elementos.listaAnimales, "Todavía no hay animales.", "empty-state empty-state--small");
    return;
  }

  const fragmento = document.createDocumentFragment();

  for (const animal of estado.animales) {
    fragmento.append(crearItemAnimal(animal));
  }

  elementos.listaAnimales.replaceChildren(fragmento);
}

function crearItemAnimal(animal) {
  const dueno = estado.usuarios.find((usuario) => usuario.id === animal.id_dueno);

  const item = document.createElement("div");
  item.className = "mini-item";

  const info = document.createElement("div");
  info.className = "mini-item__info";

  const titulo = document.createElement("span");
  titulo.className = "mini-item__title";
  titulo.textContent = `${animal.nombre} (${animal.raza})`;
  info.append(titulo);

  const detalle = document.createElement("span");
  detalle.className = "mini-item__detail";
  detalle.textContent = `Dueño: ${dueno ? nombreUsuario(dueno) : `Usuario #${animal.id_dueno}`} · ${formatearPrecio(animal.precio)}`;
  info.append(detalle);

  item.append(info, crearBotonEliminar("animal", animal.id, "Eliminar animal"));
  return item;
}

function crearDetalle(texto) {
  const detalle = document.createElement("li");
  detalle.textContent = texto;
  return detalle;
}

function crearBotonEliminar(tipo, id, etiqueta) {
  const boton = document.createElement("button");
  boton.type = "button";
  boton.className = "btn btn-danger";
  boton.dataset.action = "delete";
  boton.dataset.tipo = tipo;
  boton.dataset.id = String(id);
  boton.textContent = etiqueta;
  boton.setAttribute("aria-label", `${etiqueta} ${id}`);
  return boton;
}

function actualizarContadores() {
  elementos.contadorPublicaciones.textContent = String(estado.publicaciones.length);
  elementos.contadorUsuarios.textContent = String(estado.usuarios.length);
  elementos.contadorAnimales.textContent = String(estado.animales.length);
}

function actualizarDisponibilidadFormularios() {
  const hayUsuarios = estado.usuarios.length > 0;
  const hayAnimales = estado.animales.length > 0;

  if (elementos.publicacionSubmit.dataset.cargando !== "true") {
    elementos.publicacionSubmit.disabled = !hayUsuarios || !hayAnimales;
  }

  if (elementos.animalSubmit.dataset.cargando !== "true") {
    elementos.animalSubmit.disabled = !hayUsuarios;
  }
}

function nombreUsuario(usuario) {
  return `${usuario.nombre} ${usuario.apellido}`.trim();
}

function contactoUsuario(usuario) {
  if (!usuario.contacto) {
    return "Sin contacto";
  }

  if (typeof usuario.contacto === "string") {
    return usuario.contacto || "Sin contacto";
  }

  if (usuario.contacto.Valid) {
    return usuario.contacto.String || "Sin contacto";
  }

  return "Sin contacto";
}

function formatearPrecio(precio) {
  const valor = Number.parseFloat(precio);

  if (!Number.isFinite(valor)) {
    return precio || "Sin precio";
  }

  return new Intl.NumberFormat("es-AR", {
    style: "currency",
    currency: "ARS",
    maximumFractionDigits: 2,
  }).format(valor);
}

function mostrarMensaje(contenedor, mensaje, clase = "empty-state") {
  const parrafo = document.createElement("p");
  parrafo.className = clase;
  parrafo.textContent = mensaje;
  contenedor.replaceChildren(parrafo);
}

function mostrarStatus(elemento, mensaje, tipo = "") {
  elemento.textContent = mensaje;
  elemento.className = tipo ? `status status-${tipo}` : "status";
}

function limpiarStatus(elemento) {
  mostrarStatus(elemento, "");
}

function setBotonCargando(boton, cargando, texto = "Guardando...") {
  if (cargando) {
    boton.dataset.textoOriginal = boton.textContent;
    boton.dataset.cargando = "true";
    boton.textContent = texto;
    boton.disabled = true;
    return;
  }

  boton.textContent = boton.dataset.textoOriginal || boton.textContent;
  delete boton.dataset.textoOriginal;
  delete boton.dataset.cargando;
  boton.disabled = false;
}

async function manejarSubmitUsuario(event) {
  event.preventDefault();

  const nombre = elementos.usuarioNombre.value.trim();
  const apellido = elementos.usuarioApellido.value.trim();
  const contacto = elementos.usuarioContacto.value.trim();

  if (!nombre || !apellido) {
    mostrarStatus(elementos.usuarioStatus, "Completá nombre y apellido.", "error");
    return;
  }

  setBotonCargando(elementos.usuarioSubmit, true, "Creando...");
  limpiarStatus(elementos.usuarioStatus);

  try {
    await apiRequest(API.usuarios, {
      method: "POST",
      body: JSON.stringify({ nombre, apellido, contacto: contacto || null }),
    });

    elementos.formUsuario.reset();
    mostrarStatus(elementos.usuarioStatus, "Usuario creado correctamente.", "success");
    await cargarTodo({ mostrarCarga: false });
  } catch (error) {
    mostrarStatus(elementos.usuarioStatus, error.message, "error");
  } finally {
    setBotonCargando(elementos.usuarioSubmit, false);
    actualizarDisponibilidadFormularios();
  }
}

async function manejarSubmitAnimal(event) {
  event.preventDefault();

  const nombre = elementos.animalNombre.value.trim();
  const raza = elementos.animalRaza.value.trim();
  const idDueno = Number(elementos.animalDueno.value);
  const precio = elementos.animalPrecio.value;

  if (!nombre || !raza) {
    mostrarStatus(elementos.animalStatus, "Completá nombre y raza.", "error");
    return;
  }

  if (!idDueno) {
    mostrarStatus(elementos.animalStatus, "Seleccioná un dueño.", "error");
    return;
  }

  if (!precio || Number.parseFloat(precio) < 0) {
    mostrarStatus(elementos.animalStatus, "Ingresá un precio válido.", "error");
    return;
  }

  setBotonCargando(elementos.animalSubmit, true, "Creando...");
  limpiarStatus(elementos.animalStatus);

  try {
    await apiRequest(API.animales, {
      method: "POST",
      body: JSON.stringify({ nombre, raza, id_dueno: idDueno, precio }),
    });

    elementos.formAnimal.reset();
    mostrarStatus(elementos.animalStatus, "Animal creado correctamente.", "success");
    await cargarTodo({ mostrarCarga: false });
  } catch (error) {
    mostrarStatus(elementos.animalStatus, error.message, "error");
  } finally {
    setBotonCargando(elementos.animalSubmit, false);
    actualizarDisponibilidadFormularios();
  }
}

async function manejarSubmitPublicacion(event) {
  event.preventDefault();

  const idAnimal = Number(elementos.publicacionAnimal.value);
  const idVendedor = Number(elementos.publicacionVendedor.value);
  const precio = elementos.publicacionPrecio.value;

  if (!idAnimal) {
    mostrarStatus(elementos.publicacionStatus, "Seleccioná un animal.", "error");
    return;
  }

  if (!idVendedor) {
    mostrarStatus(elementos.publicacionStatus, "Seleccioná un vendedor.", "error");
    return;
  }

  if (!precio || Number.parseFloat(precio) < 0) {
    mostrarStatus(elementos.publicacionStatus, "Ingresá un precio válido.", "error");
    return;
  }

  setBotonCargando(elementos.publicacionSubmit, true, "Publicando...");
  limpiarStatus(elementos.publicacionStatus);

  try {
    await apiRequest(API.publicaciones, {
      method: "POST",
      body: JSON.stringify({ id_animal: idAnimal, id_vendedor: idVendedor, precio }),
    });

    elementos.formPublicacion.reset();
    mostrarStatus(elementos.publicacionStatus, "Publicación creada correctamente.", "success");
    await cargarTodo({ mostrarCarga: false });
  } catch (error) {
    mostrarStatus(elementos.publicacionStatus, error.message, "error");
  } finally {
    setBotonCargando(elementos.publicacionSubmit, false);
    actualizarDisponibilidadFormularios();
  }
}

async function manejarClickLista(event) {
  const boton = event.target.closest('button[data-action="delete"]');

  if (!boton) {
    return;
  }

  event.preventDefault();

  const tipo = boton.dataset.tipo;
  const id = Number(boton.dataset.id);
  const recurso = RECURSOS[tipo];

  if (!recurso || !Number.isInteger(id) || id <= 0) {
    return;
  }

  const confirmado = window.confirm(`¿Eliminar ${recurso.nombre} #${id}?`);

  if (!confirmado) {
    return;
  }

  const textoOriginal = boton.textContent;
  boton.disabled = true;
  boton.textContent = "Eliminando...";

  try {
    await apiRequest(recurso.ruta(id), { method: "DELETE" });
    mostrarStatus(recurso.estado(), recurso.exito, "success");
    await cargarTodo({ mostrarCarga: false });
  } catch (error) {
    mostrarStatus(recurso.estado(), error.message, "error");
  } finally {
    boton.disabled = false;
    boton.textContent = textoOriginal;
  }
}

function inicializar() {
  inicializarElementos();

  elementos.formUsuario.addEventListener("submit", manejarSubmitUsuario);
  elementos.formAnimal.addEventListener("submit", manejarSubmitAnimal);
  elementos.formPublicacion.addEventListener("submit", manejarSubmitPublicacion);
  elementos.publicacionAnimal.addEventListener("change", sincronizarVendedorConAnimal);
  document.addEventListener("click", manejarClickLista);

  cargarTodo();
}

inicializar();
