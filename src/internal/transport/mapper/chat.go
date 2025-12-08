package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/transport/dto"
)

func NewChat(chat *entity.Chat) *dto.Chat {
	return &dto.Chat{
		ID:         chat.ID,
		Type:       chat.Type,
		Name:       chat.Name,
		AvatarPath: chat.AvatarPath,
	}
}
