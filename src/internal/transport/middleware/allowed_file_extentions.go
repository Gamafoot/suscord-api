package middleware

import (
	"net/http"
	"path/filepath"
	"strings"

	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (mw *Middleware) AllowedFileExtentions() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method

			if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
				contentType := c.Request().Header.Get("Content-Type")
				if strings.Contains(contentType, "multipart/form-data") {
					file, err := c.FormFile("file")
					if err != nil {
						return utils.NewErrorResponse(c, http.StatusBadRequest, domainErrors.ErrInvalidFile.Error())
					}

					fileExt := filepath.Ext(file.Filename)
					ok := false

					for _, allowedExt := range mw.config.Media.AllowedMedia {
						if fileExt == allowedExt {
							ok = true
							break
						}
					}
					if !ok {
						return utils.NewErrorResponse(c, http.StatusBadRequest, domainErrors.ErrInvalidFileExtention.Error())
					}
				}
			}

			return next(c)
		}
	}
}
