package mapper

import (
	"suscord/internal/domain/eventbus/dto"
)

func NewMessageDelete(exceptUserID, chatID, messageID uint) *dto.MessageDelete {
	return &dto.MessageDelete{
		MessageID:    messageID,
		ChatID:       chatID,
		ExceptUserID: exceptUserID,
	}
}
