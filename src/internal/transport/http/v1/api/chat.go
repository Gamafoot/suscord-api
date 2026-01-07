package api

import (
	"context"
	"errors"
	"net/http"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/mapper"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitChatRoutes(route *echo.Group) {
	route.GET("/chats", h.GetUserChats)
	route.POST("/chats/private", h.GetOrCreatePrivateChat)
	route.POST("/chats/group", h.CreateGroupChat)
	route.PATCH("/chats/:chat_id", h.UpdateGroupChat)
	route.DELETE("/chats/:chat_id", h.DeletePrivateChat)
}

func (h *handler) GetUserChats(c echo.Context) error {
	ctx, canсel := context.WithTimeout(c.Request().Context(), h.cfg.Server.Timeout)
	defer canсel()

	userID := c.Get("user_id").(uint)

	searchParam := c.QueryParam("search")

	chats, err := h.service.Chat().GetUserChats(ctx, userID, searchParam)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result := make([]*dto.Chat, len(chats))

	for i, chat := range chats {
		result[i] = mapper.NewChat(chat, h.cfg.Media.Url)
	}

	return c.JSON(http.StatusOK, result)
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
			FriendID: reqInput.UserID,
		},
	)
	if err != nil {
		return err
	}

	result := mapper.NewChat(chat, h.cfg.Media.Url)

	return c.JSON(http.StatusOK, result)
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

	data := &entity.CreateGroupChat{Name: reqInput.Name}

	chat, err := h.service.Chat().CreateGroupChat(c.Request().Context(), userID, data)
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

	file, _ := c.FormFile("file")

	if reqInput.Name == nil && file == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(uint)

	isMember, err := h.service.ChatMember().IsMemberOfChat(c.Request().Context(), userID, chatID)
	if err != nil {
		return err
	}
	if !isMember {
		return utils.NewErrorResponse(c, http.StatusForbidden, domainErrors.ErrUserIsNotMemberOfChat.Error())
	}

	var avatarPath *string
	if file != nil {
		if !utils.FilenameValidate(file.Filename, h.cfg.Media.AllowedMedia) {
			return utils.NewErrorResponse(c, http.StatusBadRequest, domainErrors.ErrInvalidFile.Error())
		}

		if !utils.IsImage(file.Filename) {
			return utils.NewErrorResponse(c, http.StatusBadRequest, domainErrors.ErrIsNotImage.Error())
		}

		path, err := h.service.File().UploadFile(file, "chats/avatars")
		if err != nil {
			return err
		}
		avatarPath = &path
	}

	data := &entity.UpdateChat{
		Name:       reqInput.Name,
		AvatarPath: avatarPath,
	}

	chat, err := h.service.Chat().UpdateGroupChat(
		c.Request().Context(),
		userID,
		chatID,
		data,
	)
	if err != nil {
		return err
	}

	result := mapper.NewChat(chat, h.cfg.Media.Url)

	return c.JSON(http.StatusOK, result)
}

func (h *handler) DeletePrivateChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.Chat().DeletePrivateChat(c.Request().Context(), userID, chatID)
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
