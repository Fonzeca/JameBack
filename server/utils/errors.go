package utils

import "net/http"

var (
	ErrUserNotFound = NewHTTPError(http.StatusNotFound, "1-UNF", "User not found")

	ErrInternalError = NewHTTPError(http.StatusInternalServerError, "2-GEM", "Internal server error")

	ErrTryLogin = NewHTTPError(http.StatusBadRequest, "3-UNA", "Username or password incorrect")

	ErrBadRequest = NewHTTPError(http.StatusBadRequest, "4-BR", "Bad request")

	ErrUnauthorized = NewHTTPError(http.StatusUnauthorized, "5-UN", "Unauthorized")

	ErrExpiredToken = NewHTTPError(http.StatusUnauthorized, "6-ET", "Expired token")

	ErrOnInsertNoUsername = NewHTTPError(http.StatusBadRequest, "8-UNE", "Username empty")

	ErrOnInsertNoPassword = NewHTTPError(http.StatusBadRequest, "9-PE", "Password empty")

	ErrOnInsertNoDocument = NewHTTPError(http.StatusBadRequest, "10-DTE", "Document type empty")

	ErrOnChangePassword = NewHTTPError(http.StatusConflict, "11-IT", "Incorrect token")

	ErrNoBearerToken = NewHTTPError(http.StatusBadRequest, "12-TNF", "Token not found")

	ErrSamePassword = NewHTTPError(http.StatusConflict, "13-SP", "La nueva contraseña es igual a la contraseña actual")

	ErrHasChangedPassword = NewHTTPError(http.StatusBadRequest, "14-SP", "El usuario ya tuve su primer cambio de contraseña")
)

func ErrNoValidRole(roleName string) error {
	return NewHTTPError(http.StatusBadRequest, "7-NVR", "No valid role: "+roleName)
}
