package app

import (
	"log"
	"suscord/internal/config"
	"suscord/internal/database/gorm"
	"suscord/internal/eventbus"
	"suscord/internal/service"

	"golang.org/x/sync/errgroup"
)

type App struct {
	config          *config.Config
	httpServer      *httpServer
	websocketServer *websocketServer
}

func NewApp() (*App, error) {
	app := &App{
		config: config.GetConfig(),
	}

	storage, err := gorm.NewGormStorage(app.config.Database.URL)
	if err != nil {
		panic(err)
	}

	eventbus := eventbus.NewBus()

	service := service.NewService(app.config, storage, eventbus)

	app.httpServer = NewHttpServer(
		app.config,
		service,
		storage,
		eventbus,
	)

	app.websocketServer = NewWebsocketServer(
		app.config,
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
