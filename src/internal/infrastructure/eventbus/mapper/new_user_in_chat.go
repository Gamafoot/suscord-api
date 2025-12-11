package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewUserInChat(chatID uint, user *entity.User, mediaURL string) *dto.NewUserInChat {
	return &dto.NewUserInChat{
		ChatID: chatID,
		User:   NewUser(user, mediaURL),
	}
}
