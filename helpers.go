package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func parseID(valor string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(valor), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("ID inválido")
	}
	return id, nil
}

func parsePrecio(valor string) (string, error) {
	numero, err := strconv.ParseFloat(strings.TrimSpace(valor), 64)
	if err != nil || numero < 0 {
		return "", fmt.Errorf("precio inválido")
	}
	return strconv.FormatFloat(numero, 'f', 2, 64), nil
}

func convertirContacto(valor string) sql.NullString {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: valor, Valid: true}
}

func esErrorClaveForanea(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
}
