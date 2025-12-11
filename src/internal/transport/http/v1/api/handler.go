package api

import (
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/service"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/middleware"

	"github.com/labstack/echo/v4"
)

type handler struct {
	cfg        *config.Config
	service    service.Service
	storage    storage.Storage
	eventbus   eventbus.Bus
	middleware *middleware.Middleware
}

func NewHandler(
	config *config.Config,
	service service.Service,
	storage storage.Storage,
	eventbus eventbus.Bus,
	middleware *middleware.Middleware,
) *handler {
	return &handler{
		cfg:        config,
		service:    service,
		storage:    storage,
		eventbus:   eventbus,
		middleware: middleware,
	}
}

func (h *handler) InitRoutes(route *echo.Group) {
	requiredAuth := route.Group("", h.middleware.RequiredAuth())
	h.InitUserRoutes(requiredAuth)
	h.InitChatRoutes(requiredAuth)
	h.InitChatMemberRoutes(requiredAuth)
	h.InitMessageRoutes(requiredAuth)
	h.InitAttachmentRoutes(requiredAuth)
}
