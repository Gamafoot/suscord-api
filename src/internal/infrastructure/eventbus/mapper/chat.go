package mapper

import (
	"suscord/internal/domain/entity"
	"suscord/internal/domain/eventbus/dto"
	"suscord/pkg/urlpath"
)

func NewChat(chat *entity.Chat, mediaURL string) *dto.Chat {
	return &dto.Chat{
		ID:        chat.ID,
		Name:      chat.Name,
		AvatarUrl: urlpath.GetMediaURL(mediaURL, chat.AvatarPath),
		Type:      chat.Type,
	}
}
