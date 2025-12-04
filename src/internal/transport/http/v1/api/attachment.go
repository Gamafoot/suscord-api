package api

import (
	"errors"
	"net/http"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitAttachmentRoutes(route *echo.Group) {
	route.DELETE("/attachments/:attachment_id", h.DeleteAttachment)
}

func (h *handler) DeleteAttachment(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	attachmentID, err := utils.GetUIntFromParam(c, "attachment_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.Attachment().Delete(
		c.Request().Context(),
		userID,
		attachmentID,
	)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return utils.NewErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
