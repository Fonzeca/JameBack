package utils

import "net/http"

var (
	ErrUserNotFound = NewHTTPError(http.StatusNotFound, "1-UNF", "No se encontro el usuario")

	ErrInternalError = NewHTTPError(http.StatusInternalServerError, "2-GEM", "Algo paso con el servidor")

	ErrTryLogin = NewHTTPError(http.StatusBadRequest, "3-UNA", "Email o contraseña incorrecta")

	ErrBadRequest = NewHTTPError(http.StatusBadRequest, "4-BR", "Formato de llamada incorrecto")

	ErrUnauthorized = NewHTTPError(http.StatusUnauthorized, "5-UN", "Privilegios incorrectos")

	ErrExpiredToken = NewHTTPError(http.StatusUnauthorized, "6-ET", "Sesion expirada")

	ErrOnInsertNoUsername = NewHTTPError(http.StatusBadRequest, "8-UNE", "Email vacio")

	ErrOnInsertNoPassword = NewHTTPError(http.StatusBadRequest, "9-PE", "Contraseña vacia")

	ErrOnInsertNoDocument = NewHTTPError(http.StatusBadRequest, "10-DTE", "Tipo de documento vacio")

	ErrOnChangePassword = NewHTTPError(http.StatusConflict, "11-IT", "Codigo incorrecto")

	ErrNoBearerToken = NewHTTPError(http.StatusBadRequest, "12-TNF", "Sesion malformada")

	ErrSamePassword = NewHTTPError(http.StatusConflict, "13-SP", "La nueva contraseña es igual a la contraseña actual")

	ErrHasChangedPassword = NewHTTPError(http.StatusBadRequest, "14-SP", "El usuario ya tuve su primer cambio de contraseña")
)

func ErrNoValidRole(roleName string) error {
	return NewHTTPError(http.StatusBadRequest, "7-NVR", "Rol invalido: "+roleName)
}
