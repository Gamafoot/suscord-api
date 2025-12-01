package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Printf("%+v", err)

			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
		}

		log.Printf("%s %s %d", c.Request().Method, c.Request().RequestURI, c.Response().Status)

		return nil
	}
}
