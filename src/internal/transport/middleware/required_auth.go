package middleware

import (
	"net/http"
	domainErrors "suscord/internal/domain/errors"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (mw *Middleware) RequiredAuth() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie("session")
			if err != nil {
				return c.Redirect(http.StatusFound, "/auth")
			}

			session, err := mw.storage.Database().Session().GetByUUID(ctx, cookie.Value)
			if err != nil {
				if pkgErrors.Is(err, domainErrors.ErrRecordNotFound) {
					return c.Redirect(http.StatusFound, "/auth")
				}
				return err
			}

			c.Set("user_id", session.UserID)

			return next(c)
		}
	}
}
