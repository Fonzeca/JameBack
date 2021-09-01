package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initApi() {
	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router(e)

	e.Logger.Fatal(e.Start(":1323"))
}

func router(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		connect()
		return c.String(http.StatusOK, "Hello, World!")
	})

}
