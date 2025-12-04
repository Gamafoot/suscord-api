package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"suscord/internal/config"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/service"
	"suscord/internal/domain/storage"
	v1API "suscord/internal/transport/http/v1/api"
	v1WEB "suscord/internal/transport/http/v1/web"
	customMiddleware "suscord/internal/transport/middleware"

	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type httpServer struct {
	cfg     *config.Config
	echo    *echo.Echo
	service service.Service
	storage storage.Storage
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewHttpServer(cfg *config.Config, services service.Service, storage storage.Storage, eventbus eventbus.Bus) *httpServer {
	server := &httpServer{
		cfg:     cfg,
		echo:    echo.New(),
		service: services,
		storage: storage,
	}

	server.echo.Static(server.cfg.Static.RootUrl, server.cfg.Static.RootFolder)
	server.echo.Static(server.cfg.Media.RootUrl, server.cfg.Media.RootFolder)

	template := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("assets/html/*.html")),
	}
	server.echo.Renderer = template
	server.echo.Validator = &CustomValidator{validator: validator.New()}

	_customMiddleware := customMiddleware.NewMiddleware(storage)

	// server.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: server.cfg.CORS.Origins,
	// 	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	// }))
	// server.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	// 	Timeout: server.cfg.Server.Timeout,
	// }))

	handlerV1WEB := v1WEB.NewHandler(
		server.cfg,
		server.service,
		server.storage,
		_customMiddleware,
	)

	route := server.echo.Group(
		"",
		middleware.BodyLimit(server.cfg.Media.MaxSize),
		customMiddleware.Logger,
		_customMiddleware.AllowedFileExtentions(),
	)

	handlerV1WEB.InitRoutes(route)

	handlerV1API := v1API.NewHandler(
		server.cfg,
		server.service,
		server.storage,
		eventbus,
		_customMiddleware,
	)
	handlerV1API.InitRoutes(server.echo.Group("/api/v1"))

	return server
}

func (s *httpServer) Run() error {
	host := flag.String("host", "", "Пример -host=0.0.0.0")
	addr := fmt.Sprintf("%s:%s", *host, s.cfg.Server.Port)

	if err := s.echo.Start(addr); err != nil {
		return err
	}

	return nil
}

func (s *httpServer) ShutdownServer(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *httpServer) Echo() *echo.Echo {
	return s.echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
