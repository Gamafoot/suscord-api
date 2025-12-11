package api

import (
	"errors"
	"net/http"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/mapper"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitChatMemberRoutes(route *echo.Group) {
	route.GET("/chats/:chat_id/members", h.GetChatMembers)
	route.GET("/chats/:chat_id/non-members", h.GetNonMembers)
	route.POST("/chats/:chat_id/invite", h.SendInviteChat)
	route.GET("/chats/invite/accept/:code", h.AcceptInviteChat)
	route.GET("/chats/:chat_id/leave", h.LeaveFromChat)
}

func (h *handler) GetChatMembers(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	users, err := h.service.ChatMember().GetNonMembers(c.Request().Context(), chatID, userID)
	if err != nil {
		return err
	}

	result := make([]*dto.User, len(users))
	for i, user := range users {
		result[i] = mapper.NewUser(user, h.cfg.Media.Url)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *handler) GetNonMembers(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	users, err := h.service.ChatMember().GetNotChatMembers(c.Request().Context(), chatID, userID)
	if err != nil {
		return err
	}

	result := make([]*dto.User, len(users))
	for i, user := range users {
		result[i] = mapper.NewUser(user, h.cfg.Media.Url)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *handler) SendInviteChat(c echo.Context) error {
	reqInput := new(dto.InviteUserRequest)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.ChatMember().SendInvite(
		c.Request().Context(),
		userID,
		chatID,
		reqInput.UserID,
	)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) AcceptInviteChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	code := c.Param("code")

	err := h.service.ChatMember().AcceptInvite(c.Request().Context(), userID, code)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRedisNil) {
			return utils.NewErrorResponse(c, http.StatusNotFound, "wrong invite code")
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) LeaveFromChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.ChatMember().LeaveFromChat(c.Request().Context(), chatID, userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserIsNotMemberOfChat) {
			return utils.NewErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
