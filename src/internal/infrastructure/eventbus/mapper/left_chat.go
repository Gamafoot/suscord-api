package mapper

import "suscord/internal/domain/eventbus/dto"

func NewLeftChat(chatID, userID uint) *dto.LeftChat {
	return &dto.LeftChat{
		ChatID: chatID,
		UserID: userID,
	}
}
