package middleware

import (
	"net/http"

	domainErrors "suscord/internal/domain/errors"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (mw *Middleware) NotAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, _ := c.Cookie("session")
			if cookie != nil {
				session, err := mw.storage.Database().Session().GetByUUID(c.Request().Context(), cookie.Value)
				if err != nil {
					if !pkgErrors.Is(err, domainErrors.ErrRecordNotFound) {
						return err
					}
				}

				if session != nil {
					return c.Redirect(http.StatusSeeOther, "/")
				}
			}

			return next(c)
		}
	}
}
