package main

import (
	"net/http"
	db "tp2/db/sqlc"
)

// esta estrucuta se hizo para poder comunicar y poder usar los metodos que nos creo sqlc
type estructuraUsuario struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}

func (h *estructuraUsuario) CRUDusuarios(w http.ResponseWriter, r *http.Request) {
}

func (h *estructuraUsuario) GetUsuarioByID(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) GetUsuarioByName(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) GetUsuarioList(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) CreateUsuario(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) UpdateUsuario(w http.ResponseWriter, r *http.Request) {

}

func (h *estructuraUsuario) DeleteUsuario(w http.ResponseWriter, r *http.Request) {

}
