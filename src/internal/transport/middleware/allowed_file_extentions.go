package middleware

import (
	"net/http"
	"path/filepath"

	domainErrors "suscord/internal/domain/errors"

	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

func (mw *Middleware) AllowedFileExtentions() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method

			if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
				contentType := c.Request().Header.Get("Content-Type")
				if contentType == "multipart/form-data" {
					file, err := c.FormFile("file")
					if err != nil {
						return pkgErrors.WithStack(domainErrors.ErrInvalidFile)
					}

					fileExt := filepath.Ext(file.Filename)
					ok := false

					for _, allowedExt := range mw.config.Media.AllowedExtentions {
						if fileExt == allowedExt {
							ok = true
							break
						}
					}
					if !ok {
						return pkgErrors.WithStack(domainErrors.ErrInvalidFile)
					}
				}
			}

			return next(c)
		}
	}
}
