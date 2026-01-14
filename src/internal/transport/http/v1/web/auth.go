package web

import (
	"net/http"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (h *handler) InitAuthRoutes(route *echo.Group) {
	route.GET("/auth", h.AuthPage)
	route.POST("/auth", h.AuthPage)
}

func (h *handler) AuthPage(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		uuid, err := h.login(c)
		if err != nil {
			data := make(map[string]string, 0)
			if pkgErrors.Is(err, domainErrors.ErrInvalidLoginOrPassword) {
				data["error"] = "Неверный логин или пароль. Попробуйте еще раз."
				return c.Render(http.StatusOK, "auth.html", data)
			}
			return err
		}

		c.SetCookie(&http.Cookie{
			Name:  "session",
			Value: uuid,
		})
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "auth.html", nil)
}

func (h *handler) login(c echo.Context) (string, error) {
	input := &dto.LoginOrCreateRequest{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return "", err
	}

	uuid, err := h.service.Auth().Login(c.Request().Context(), &entity.LoginOrCreateInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (h *handler) GetSession(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	userID := sess.Values["user_id"].(int)

	return c.JSON(http.StatusOK, map[string]int{
		"user_id": userID,
	})
}
