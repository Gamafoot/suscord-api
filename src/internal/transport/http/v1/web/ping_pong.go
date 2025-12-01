package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitPingPongRoutes(route *echo.Group) {
	route.POST("/ping", h.PingPong)
}

func (h *handler) PingPong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
