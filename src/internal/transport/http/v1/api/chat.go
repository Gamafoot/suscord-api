package api

import (
	"context"
	"errors"
	"net/http"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitChatRoutes(route *echo.Group) {
	route.GET("/chats", h.GetUserChats)
	route.POST("/chats/private", h.GetOrCreatePrivateChat)
	route.POST("/chats", h.CreateGroupChat)
	route.PATCH("/chats/:chat_id", h.UpdateGroupChat)
	route.DELETE("/chats/:chat_id", h.DeletePrivateChat)
}

func (h *handler) GetUserChats(c echo.Context) error {
	ctx, canсel := context.WithTimeout(c.Request().Context(), h.config.Server.Timeout)
	defer canсel()

	userID := c.Get("user_id").(uint)

	chats, err := h.service.Chat().GetUserChats(ctx, userID)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	respChats := make([]*dto.Chat, len(chats))

	for i, chat := range chats {
		respChats[i] = &dto.Chat{
			ID:         chat.ID,
			Type:       chat.Type,
			Name:       chat.Name,
			AvatarPath: chat.AvatarPath,
		}
	}

	return c.JSON(http.StatusOK, respChats)
}

func (h *handler) GetOrCreatePrivateChat(c echo.Context) error {
	reqInput := new(dto.CreatePrivateChatRequest)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userID := c.Get("user_id").(uint)

	chat, err := h.service.Chat().GetOrCreatePrivateChat(
		c.Request().Context(),
		&entity.CreatePrivateChat{
			UserID:   userID,
			FriendID: reqInput.FriendID,
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, chat)
}

func (h *handler) CreateGroupChat(c echo.Context) error {
	reqInput := new(dto.CreateGroupChatRequest)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userID := c.Get("user_id").(uint)

	chat, err := h.service.Chat().CreateGroupChat(c.Request().Context(), &entity.CreateGroupChat{
		UserID:     userID,
		FriendID:   reqInput.FriendID,
		Name:       reqInput.Name,
		AvatarPath: reqInput.AvatarPath,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, chat)
}

func (h *handler) UpdateGroupChat(c echo.Context) error {
	reqInput := new(dto.UpdateGroupChatInput)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if reqInput.Name == nil && reqInput.AvatarPath == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(uint)

	chat, err := h.service.Chat().UpdateGroupChat(
		c.Request().Context(),
		userID,
		chatID,
		&entity.UpdateChat{
			Name:       reqInput.Name,
			AvatarPath: reqInput.AvatarPath,
		},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, chat)
}

func (h *handler) DeletePrivateChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.Chat().DeletePrivateChat(c.Request().Context(), chatID, userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrForbidden) {
			return utils.NewErrorResponse(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return utils.NewErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
