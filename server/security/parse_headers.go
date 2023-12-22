package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func ParseHeadersMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Leer los headers
		headers := c.Request().Header

		// Imprimir todos los headers
		for key, values := range headers {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		// Verificar si el header "X-Username" existe
		if headers.Get("X-Username") == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized by ParseHeadersMiddleware")
		}

		// Pasar los headers al contexto de Echo
		roles := strings.Split(headers.Get("X-Roles"), ",")
		c.Set("username", headers.Get("X-Username"))
		c.Set("admin", headers.Get("X-Admin"))
		c.Set("roles", roles)

		// Llamar al siguiente handler
		return next(c)
	}
}
