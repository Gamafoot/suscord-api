package app

import (
	"log"
	"suscord/internal/config"
	"suscord/internal/infrastructure/eventbus"
	"suscord/internal/infrastructure/service"

	"golang.org/x/sync/errgroup"
)

type App struct {
	httpServer      *httpServer
	websocketServer *websocketServer
}

func NewApp() (*App, error) {
	cfg := config.GetConfig()

	storage, err := NewStorage(cfg)
	if err != nil {
		panic(err)
	}

	eventbus := eventbus.NewBus()

	service := service.NewService(cfg, storage, eventbus)

	app := new(App)

	app.httpServer = NewHttpServer(
		cfg,
		service,
		storage,
		eventbus,
	)

	app.websocketServer = NewWebsocketServer(
		cfg,
		app.httpServer.Echo(),
		storage,
		eventbus,
	)

	return app, nil
}

func (a *App) Run() error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		log.Println("Websocker server is running")
		return a.websocketServer.Run()
	})

	eg.Go(func() error {
		return a.httpServer.Run()
	})

	return eg.Wait()
}
