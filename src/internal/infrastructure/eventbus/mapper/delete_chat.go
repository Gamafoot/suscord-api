package mapper

import "suscord/internal/domain/eventbus/dto"

func NewDeleteChat(chatID, exceptUserID uint) *dto.DeleteChat {
	return &dto.DeleteChat{
		ChatID:       chatID,
		ExceptUserID: exceptUserID,
	}
}
