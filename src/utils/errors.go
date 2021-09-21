package utils

import "net/http"

var (
	ErrUserNotFound = NewHTTPError(http.StatusNotFound, "1-UNF", "User not found")

	ErrInternalError = NewHTTPError(http.StatusInternalServerError, "2-GEM", "Internal server error")

	ErrTryLogin = NewHTTPError(http.StatusBadRequest, "3-UNA", "Username or password incorrect")

	ErrBadRequestGetuser = NewHTTPError(http.StatusBadRequest, "4-BR", "Bad request")

	ErrUnauthorized = NewHTTPError(http.StatusUnauthorized, "5-UN", "Unauthorized")

	ErrExpiredToken = NewHTTPError(http.StatusUnauthorized, "6-ET", "Expired token")
)
