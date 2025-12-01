package utils

import (
	"github.com/labstack/echo/v4"
)

func NewErrorResponse(c echo.Context, status int, message string) error {
	return c.String(status, message)
}
