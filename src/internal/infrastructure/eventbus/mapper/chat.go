package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
)

func NewChat(chat *entity.Chat) *dto.Chat {
	return &dto.Chat{
		ID:         chat.ID,
		Name:       chat.Name,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
	}
}
