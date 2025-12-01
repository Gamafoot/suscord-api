package web

import (
	"suscord/internal/config"
	"suscord/internal/domain/service"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	config     *config.Config
	service    service.Service
	storage    storage.Storage
	middleware *middleware.Middleware
}

func NewHandler(
	config *config.Config,
	services service.Service,
	storage storage.Storage,
	middleware *middleware.Middleware,
) *handler {
	return &handler{
		config:     config,
		service:    services,
		storage:    storage,
		middleware: middleware,
	}
}

func (h *handler) InitRoutes(route *echo.Group) {
	requiredAuth := route.Group("", h.middleware.RequiredAuth())
	{
		h.InitChatRoutes(requiredAuth)
		h.InitPingPongRoutes(requiredAuth)
	}
	h.InitAuthRoutes(route.Group("", h.middleware.NotAuth()))
}
