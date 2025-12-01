package api

import (
	"net/http"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitChatMemberRoutes(route *echo.Group) {
	route.GET("/chats/:chat_id/members", h.GetChatMembers)
	route.POST("/chats/:chat_id/members", h.AddUserToChat)
}

func (h *handler) GetChatMembers(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	users, err := h.service.ChatMember().GetChatMembers(c.Request().Context(), chatID, userID)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *handler) AddUserToChat(c echo.Context) error {
	reqInput := new(dto.AddUserToChatInput)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	memberID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.ChatMember().AddUserToChat(c.Request().Context(), chatID, memberID, reqInput.UserID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
