package api

import (
	"errors"
	"net/http"
	"suscord/internal/domain/entity"
	domainErrors "suscord/internal/domain/errors"
	"suscord/internal/transport/dto"
	"suscord/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

func (h *handler) InitMessageRoutes(route *echo.Group) {
	route.GET("/chats/:chat_id/messages", h.GetChatMessages)
	route.POST("/chats/:chat_id/messages", h.CreateMessage)
	route.PATCH("/messages/:message_id", h.UpdateMessage)
	route.DELETE("/messages/:message_id", h.DeleteMessage)
}

func (h *handler) GetChatMessages(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chatID, err := utils.GetUIntFromParam(c, "chat_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	lastMessageID, err := utils.GetIntFromQuery(c, "last_message_id", 0)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	limit, err := utils.GetIntFromQuery(c, "limit", 10)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	messages, err := h.service.Message().GetChatMessages(
		c.Request().Context(),
		&entity.GetMessagesInput{
			ChatID:        chatID,
			UserID:        userID,
			LastMessageID: uint(lastMessageID),
			Limit:         limit,
		},
	)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result := make([]*dto.MessageResponse, len(messages))

	for i, message := range messages {
		result[i] = dto.NewMessageReponse(message)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *handler) CreateMessage(c echo.Context) error {
	reqInput := new(dto.CreateMessageRequest)

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

	form, err := c.MultipartForm()
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) {
			return utils.NewErrorResponse(c, http.StatusBadRequest, "request does not contain multipart form")
		}
		return err
	}

	files := form.File["file"]

	data := &entity.CreateMessage{
		Type:    reqInput.Type,
		Content: reqInput.Content,
	}

	message, err := h.service.Message().Create(
		c.Request().Context(),
		userID,
		chatID,
		data,
		files,
	)
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	messageResp := dto.NewMessageReponse(message)

	return c.JSON(http.StatusOK, messageResp)
}

func (h *handler) UpdateMessage(c echo.Context) error {
	reqInput := new(dto.UpdateMessageRequest)

	if err := c.Bind(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(reqInput); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userID := c.Get("user_id").(uint)

	messageID, err := utils.GetUIntFromParam(c, "message_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	data := &entity.UpdateMessage{
		Content: reqInput.Content,
	}

	message, err := h.service.Message().Update(
		c.Request().Context(),
		userID,
		messageID,
		data,
	)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return c.JSON(http.StatusOK, message)
}

func (h *handler) DeleteMessage(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	messageID, err := utils.GetUIntFromParam(c, "message_id")
	if err != nil {
		return utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.service.Message().Delete(c.Request().Context(), userID, messageID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return err
	}

	return c.NoContent(http.StatusOK)
}
