package app

import (
	"log"
	"suscord/internal/config"
	"suscord/internal/database/relational"
	"suscord/internal/eventbus"
	"suscord/internal/service"

	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg             *config.Config
	httpServer      *httpServer
	websocketServer *websocketServer
}

func NewApp() (*App, error) {
	app := &App{
		cfg: config.GetConfig(),
	}

	storage, err := relational.NewGormStorage(app.cfg.Database.URL, app.cfg.Database.LogLevel)
	if err != nil {
		panic(err)
	}

	eventbus := eventbus.NewBus()

	service := service.NewService(app.cfg, storage, eventbus)

	app.httpServer = NewHttpServer(
		app.cfg,
		service,
		storage,
		eventbus,
	)

	app.websocketServer = NewWebsocketServer(
		app.cfg,
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
