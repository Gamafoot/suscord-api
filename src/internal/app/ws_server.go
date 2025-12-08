package app

import (
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/ws"
	"suscord/internal/transport/ws/hub"

	"github.com/labstack/echo/v4"
)

type websocketServer struct {
	echo     *echo.Echo
	eventbus eventbus.Bus
	hub      *hub.Hub
}

func NewWebsocketServer(
	cfg *config.Config,
	echo *echo.Echo,
	storage storage.Storage,
	eventbus eventbus.Bus,
) *websocketServer {
	hubInstance := hub.NewHub(cfg, storage, eventbus)

	server := &websocketServer{
		echo:     echo,
		eventbus: eventbus,
		hub:      hubInstance,
	}

	handler := ws.NewHandler(hubInstance)
	handler.InitRoutes(server.echo)

	return server
}

func (s *websocketServer) Run() error {
	s.hub.Run()
	return nil
}
