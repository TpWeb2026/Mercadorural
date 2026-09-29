package main


import (
	db "tp2/db/sqlc"
)

// esta estrucuta se hizo para poder comunicar y poder usar los metodos que nos creo sqlc
type estructuraUsuario struct {
	//en el main en nuestro caso, lo llamamos metodoBD, que seria nuestro *db.Queries
	Consultas *db.Queries
}