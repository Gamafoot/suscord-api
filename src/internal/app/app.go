package app

import (
	"suscord/internal/config"
	"suscord/internal/infrastructure/broker/rabbitmq"
	"suscord/internal/infrastructure/service"
	"suscord/pkg/logger"
)

type App struct {
	httpServer *httpServer
}

func NewApp() (*App, error) {
	cfg := config.GetConfig()

	storage, err := NewStorage(cfg)
	if err != nil {
		return nil, err
	}

	logger, err := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Folder)
	if err != nil {
		return nil, err
	}

	broker, err := rabbitmq.NewBroker(cfg.Broker.Addr, cfg.Broker.PoolSize, logger)
	if err != nil {
		return nil, err
	}

	service := service.NewService(cfg, storage, broker, logger)

	app := new(App)
	app.httpServer = NewHttpServer(cfg, service, storage)

	return app, nil
}

func (a *App) Run() error {
	return a.httpServer.Run()
}
