package app

import (
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/ws"

	"github.com/labstack/echo/v4"
)

type websocketServer struct {
	echo     *echo.Echo
	eventbus eventbus.Bus
}

func NewWebsocketServer(cfg *config.Config, echo *echo.Echo, storage storage.Storage, eventbus eventbus.Bus) *websocketServer {
	s := &websocketServer{
		echo:     echo,
		eventbus: eventbus,
	}

	handler := ws.NewHandler(cfg, storage, s.eventbus)
	handler.InitRoutes(s.echo)

	return s
}

func (s *websocketServer) Run() error {
	return nil
}
