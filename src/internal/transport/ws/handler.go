package ws

import (
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/ws/hub"

	"github.com/labstack/echo/v4"
)

type handler struct {
	storage storage.Storage
	hub     *hub.Hub
}

func NewHandler(cfg *config.Config, storage storage.Storage, eventbus eventbus.Bus) *handler {
	hub := hub.NewHub(cfg, storage, eventbus)
	go hub.Run()
	return &handler{
		hub: hub,
	}
}

func (h *handler) InitRoutes(route *echo.Echo) {
	route.GET("/ws", h.hub.WebsocketHandler)
}
